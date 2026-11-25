package service

import (
	"archive/tar"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"golang.org/x/crypto/argon2"
)

var (
	storageKeyMu sync.RWMutex
	storageKey   []byte
)

// Magic bytes for Aparok backup container
var backupMagic = []byte("APAROK_BACKUP_V1\n")

// SetStorageKey supplies the 32-byte encryption key from Android Keystore.
func SetStorageKey(key []byte) error {
	if len(key) != 32 {
		return errors.New("storage key must be exactly 32 bytes")
	}
	storageKeyMu.Lock()
	defer storageKeyMu.Unlock()

	// Zeroize old key if any
	if storageKey != nil {
		for i := range storageKey {
			storageKey[i] = 0
		}
	}
	storageKey = make([]byte, 32)
	copy(storageKey, key)
	return nil
}

// Lock zeroizes the in-memory storage key and locks the core.
func Lock() {
	storageKeyMu.Lock()
	defer storageKeyMu.Unlock()
	if storageKey != nil {
		for i := range storageKey {
			storageKey[i] = 0
		}
		storageKey = nil
	}
}

// IsLocked reports whether the core storage key is missing/zeroized.
func IsLocked() bool {
	storageKeyMu.RLock()
	defer storageKeyMu.RUnlock()
	return storageKey == nil
}

// GetStorageKey returns a copy of current storage key or an error if locked.
func GetStorageKey() ([]byte, error) {
	storageKeyMu.RLock()
	defer storageKeyMu.RUnlock()
	if storageKey == nil {
		return nil, errors.New("core is locked: storage key unavailable")
	}
	k := make([]byte, 32)
	copy(k, storageKey)
	return k, nil
}

// ProbeStorage tests if a directory supports full database semantics:
// write, fsync, atomic rename, file lock, and adequate free space.
func ProbeStorage(targetDir string) (string, error) {
	if targetDir == "" {
		return "", errors.New("directory path cannot be empty")
	}
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return "", fmt.Errorf("cannot create storage directory: %w", err)
	}

	testFile := filepath.Join(targetDir, fmt.Sprintf(".probe_%d.tmp", os.Getpid()))
	renamedFile := filepath.Join(targetDir, fmt.Sprintf(".probe_renamed_%d.tmp", os.Getpid()))
	defer os.Remove(testFile)
	defer os.Remove(renamedFile)

	// 1. Write test
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return "", fmt.Errorf("write test failed: %w", err)
	}

	payload := []byte("aparok-sovereign-storage-probe-data")
	if _, err := f.Write(payload); err != nil {
		f.Close()
		return "", fmt.Errorf("write payload failed: %w", err)
	}

	// 2. Fsync test
	if err := f.Sync(); err != nil {
		f.Close()
		return "", fmt.Errorf("fsync test failed: %w", err)
	}

	// 3. File lock test (flock)
	fd := int(f.Fd())
	if err := syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return "", fmt.Errorf("file lock test failed: %w", err)
	}
	_ = syscall.Flock(fd, syscall.LOCK_UN)
	f.Close()

	// 4. Atomic rename test
	if err := os.Rename(testFile, renamedFile); err != nil {
		return "", fmt.Errorf("atomic rename test failed: %w", err)
	}

	// 5. Statfs free space check (minimum 50 MB)
	var stat syscall.Statfs_t
	if err := syscall.Statfs(targetDir, &stat); err == nil {
		freeBytes := stat.Bavail * uint64(stat.Bsize)
		if freeBytes < 50*1024*1024 {
			return "", fmt.Errorf("insufficient free space: %d MB available, need at least 50 MB", freeBytes/(1024*1024))
		}
	}

	return "Storage verified: full database semantics supported", nil
}

type ManifestEntry struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}

type BackupManifest struct {
	Version   int             `json:"version"`
	CreatedAt int64           `json:"created_at"`
	Files     []ManifestEntry `json:"files"`
}

func deriveBackupKey(password string, salt []byte) []byte {
	// Argon2id with 64MB memory, 3 iterations, 4 parallelism, 32-byte key
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
}

// Backup creates a quiesced, authenticated encrypted archive of sourceDataDir.
func Backup(sourceDataDir string, destBackupFile string, password string) error {
	if len(password) < 10 {
		return errors.New("password must be at least 10 characters")
	}
	if _, err := os.Stat(sourceDataDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDataDir)
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}
	key := deriveBackupKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	var tarBuf bytes.Buffer
	tarWriter := tar.NewWriter(&tarBuf)

	var manifest BackupManifest
	manifest.Version = 1

	err = filepath.Walk(sourceDataDir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceDataDir, p)
		if err != nil {
			return nil
		}
		// Skip temporary or lock files
		base := filepath.Base(p)
		if base == "LOCK" || filepath.Ext(base) == ".tmp" || filepath.Ext(base) == ".lock" {
			return nil
		}

		data, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", p, err)
		}
		h := sha256.Sum256(data)
		hStr := fmt.Sprintf("%x", h[:])

		manifest.Files = append(manifest.Files, ManifestEntry{
			Path: rel,
			Size: info.Size(),
			Hash: hStr,
		})

		hdr := &tar.Header{
			Name: rel,
			Mode: 0600,
			Size: int64(len(data)),
		}
		if err := tarWriter.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tarWriter.Write(data); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Write manifest.json into tar
	manifestBytes, _ := json.Marshal(manifest)
	mHdr := &tar.Header{
		Name: "manifest.json",
		Mode: 0600,
		Size: int64(len(manifestBytes)),
	}
	_ = tarWriter.WriteHeader(mHdr)
	_, _ = tarWriter.Write(manifestBytes)
	tarWriter.Close()

	// Encrypt tar buffer
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := gcm.Seal(nil, nonce, tarBuf.Bytes(), backupMagic)

	// Write output file: MAGIC || SALT(16) || NONCE(12) || CIPHERTEXT
	outF, err := os.OpenFile(destBackupFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outF.Close()

	if _, err := outF.Write(backupMagic); err != nil {
		return err
	}
	if _, err := outF.Write(salt); err != nil {
		return err
	}
	if _, err := outF.Write(nonce); err != nil {
		return err
	}
	if _, err := outF.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

// VerifyBackup decrypts and verifies manifest and integrity hashes without extracting.
func VerifyBackup(backupPath string, password string) error {
	raw, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}
	magicLen := len(backupMagic)
	if len(raw) < magicLen+16+12 {
		return errors.New("backup file too small or corrupted")
	}
	if !bytes.Equal(raw[:magicLen], backupMagic) {
		return errors.New("invalid backup format magic")
	}

	salt := raw[magicLen : magicLen+16]
	nonce := raw[magicLen+16 : magicLen+28]
	ciphertext := raw[magicLen+28:]

	key := deriveBackupKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, backupMagic)
	if err != nil {
		return errors.New("decryption failed: incorrect password or corrupted backup")
	}

	// Verify tar and manifest
	tarReader := tar.NewReader(bytes.NewReader(plaintext))
	var manifest *BackupManifest
	fileMap := make(map[string][]byte)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read error: %w", err)
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, tarReader); err != nil {
			return err
		}
		if hdr.Name == "manifest.json" {
			var m BackupManifest
			if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
				return fmt.Errorf("manifest decode failed: %w", err)
			}
			manifest = &m
		} else {
			fileMap[hdr.Name] = buf.Bytes()
		}
	}

	if manifest == nil {
		return errors.New("backup contains no manifest")
	}

	for _, entry := range manifest.Files {
		content, ok := fileMap[entry.Path]
		if !ok {
			return fmt.Errorf("missing file from manifest: %s", entry.Path)
		}
		h := sha256.Sum256(content)
		if fmt.Sprintf("%x", h[:]) != entry.Hash {
			return fmt.Errorf("checksum mismatch for %s", entry.Path)
		}
	}

	return nil
}

// Restore decrypts, verifies, and stages restore before committing to destDir.
func Restore(backupPath string, password string, destDir string) error {
	raw, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}
	magicLen := len(backupMagic)
	if len(raw) < magicLen+16+12 {
		return errors.New("backup file too small or corrupted")
	}
	if !bytes.Equal(raw[:magicLen], backupMagic) {
		return errors.New("invalid backup format magic")
	}

	salt := raw[magicLen : magicLen+16]
	nonce := raw[magicLen+16 : magicLen+28]
	ciphertext := raw[magicLen+28:]

	key := deriveBackupKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, backupMagic)
	if err != nil {
		return errors.New("decryption failed: incorrect password or corrupted backup")
	}

	stagingDir := destDir + ".staging"
	_ = os.RemoveAll(stagingDir)
	if err := os.MkdirAll(stagingDir, 0700); err != nil {
		return err
	}
	defer os.RemoveAll(stagingDir)

	tarReader := tar.NewReader(bytes.NewReader(plaintext))
	var manifest *BackupManifest

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Name == "manifest.json" {
			var m BackupManifest
			if err := json.NewDecoder(tarReader).Decode(&m); err != nil {
				return err
			}
			manifest = &m
			continue
		}
		targetFile := filepath.Join(stagingDir, hdr.Name)
		_ = os.MkdirAll(filepath.Dir(targetFile), 0700)
		f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, tarReader); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	if manifest == nil {
		return errors.New("backup manifest missing")
	}

	// Verify all staged files against manifest
	for _, entry := range manifest.Files {
		stagedPath := filepath.Join(stagingDir, entry.Path)
		data, err := os.ReadFile(stagedPath)
		if err != nil {
			return fmt.Errorf("staged file missing: %s", entry.Path)
		}
		h := sha256.Sum256(data)
		if fmt.Sprintf("%x", h[:]) != entry.Hash {
			return fmt.Errorf("checksum mismatch in restored file: %s", entry.Path)
		}
	}

	// Atomic replace of destDir
	backupOld := destDir + ".old"
	_ = os.RemoveAll(backupOld)
	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, backupOld); err != nil {
			return fmt.Errorf("cannot backup existing data: %w", err)
		}
	}

	if err := os.Rename(stagingDir, destDir); err != nil {
		// Rollback if possible
		_ = os.Rename(backupOld, destDir)
		return fmt.Errorf("atomic commit failed: %w", err)
	}

	_ = os.RemoveAll(backupOld)
	return nil
}
