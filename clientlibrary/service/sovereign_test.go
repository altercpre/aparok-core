package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorageKeyAndLock(t *testing.T) {
	// Initially locked
	Lock()
	require.True(t, IsLocked())

	// Invalid key length
	err := SetStorageKey([]byte("too-short"))
	require.Error(t, err)

	// Valid 32-byte key
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	err = SetStorageKey(key)
	require.NoError(t, err)
	require.False(t, IsLocked())

	retrieved, err := GetStorageKey()
	require.NoError(t, err)
	require.Equal(t, key, retrieved)

	// Lock zeroizes
	Lock()
	require.True(t, IsLocked())
	_, err = GetStorageKey()
	require.Error(t, err)
}

func TestProbeStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "aparok_probe_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	msg, err := ProbeStorage(tempDir)
	require.NoError(t, err)
	require.Contains(t, msg, "full database semantics supported")
}

func TestBackupVerifyRestore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "aparok_backup_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source_data")
	err = os.MkdirAll(sourceDir, 0700)
	require.NoError(t, err)

	// Create test files
	file1 := filepath.Join(sourceDir, "notes.txt")
	err = os.WriteFile(file1, []byte("sovereign aparok test note"), 0600)
	require.NoError(t, err)

	subdir := filepath.Join(sourceDir, "media")
	err = os.MkdirAll(subdir, 0700)
	require.NoError(t, err)
	file2 := filepath.Join(subdir, "avatar.bin")
	err = os.WriteFile(file2, []byte("image-data-payload-12345"), 0600)
	require.NoError(t, err)

	backupFile := filepath.Join(tempDir, "aparok-backup.aparok")
	pwd := "supersecretpassphrase123"

	// 1. Create backup
	err = Backup(sourceDir, backupFile, pwd)
	require.NoError(t, err)

	// 2. Verify backup
	err = VerifyBackup(backupFile, pwd)
	require.NoError(t, err)

	// Wrong password fails
	err = VerifyBackup(backupFile, "wrongpassword123")
	require.Error(t, err)

	// 3. Restore to target dir
	destDir := filepath.Join(tempDir, "restored_data")
	err = Restore(backupFile, pwd, destDir)
	require.NoError(t, err)

	// Verify restored files exist and match
	restored1, err := os.ReadFile(filepath.Join(destDir, "notes.txt"))
	require.NoError(t, err)
	require.Equal(t, []byte("sovereign aparok test note"), restored1)

	restored2, err := os.ReadFile(filepath.Join(destDir, "media", "avatar.bin"))
	require.NoError(t, err)
	require.Equal(t, []byte("image-data-payload-12345"), restored2)
}
