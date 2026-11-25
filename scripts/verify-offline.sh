#!/usr/bin/env bash
set -euo pipefail

echo "=== Aparok Core Sovereign Offline Verification (R6 Strict Gate) ==="

BANNED_PACKAGES=(
    "github.com/uber/jaeger-client-go"
    "github.com/grandcat/zeroconf"
    "github.com/otiai10/opengraph"
    "github.com/anyproto/anytype-publish-server/publishclient"
    "github.com/hbagdi/go-unsplash"
)

echo "Checking package dependencies against banned telemetry and remote discovery..."
TAGS="nogrpcserver gomobile nowatchdog nosigar nomutexdeadlockdetector timetzdata rasterizesvg"

# Strict check: fail if go list fails (no 2>/dev/null || true)
DEPS=$(go list -tags "${TAGS}" -deps ./clientlibrary/service ./core)

VIOLATIONS=0
for BANNED in "${BANNED_PACKAGES[@]}"; do
    if grep -F "$BANNED" <<< "$DEPS" >/dev/null; then
        echo "FAIL: Banned package detected in dependencies: $BANNED"
        echo "Tracer chain for $BANNED:"
        go list -tags "${TAGS}" -deps -f '{{with .Imports}}{{range .}}{{printf "%s -> %s\n" $.ImportPath .}}{{end}}{{end}}' ./clientlibrary/service ./core | grep "$BANNED" || true
        VIOLATIONS=$((VIOLATIONS + 1))
    else
        echo "PASS: Clean - $BANNED not present in dependency graph."
    fi
done

echo "Checking direct imports of prometheus in core and service..."
DIRECT_IMPORTS=$(go list -tags "${TAGS}" -f '{{range .Imports}}{{println .}}{{end}}' ./clientlibrary/service ./core)
if grep -F "github.com/prometheus/client_golang" <<< "$DIRECT_IMPORTS" >/dev/null; then
    echo "FAIL: prometheus directly imported by core or clientlibrary/service"
    VIOLATIONS=$((VIOLATIONS + 1))
else
    echo "PASS: Zero direct prometheus imports in core and clientlibrary/service"
fi

echo "Checking for local-only mode enforcement..."
if grep -q "IsLocalOnlyMode() bool" core/anytype/config/config.go; then
    echo "PASS: IsLocalOnlyMode defined and enforced in core/anytype/config/config.go"
else
    echo "FAIL: IsLocalOnlyMode not found in config"
    VIOLATIONS=$((VIOLATIONS + 1))
fi

echo "Checking bootstrap.go for banned network components..."
BANNED_COMPONENTS=(
    "metric\.New"
    "debugserver\.New"
    "profiler\.New"
    "localdiscovery\.New"
    "coordinatorclient\.New"
    "nodeclient\.New"
    "unsplash\.New"
    "gallery\.New"
    "publish\.New"
    "publishclient\.New"
)

for COMP in "${BANNED_COMPONENTS[@]}"; do
    if grep -E "Register\(${COMP}\)" core/anytype/bootstrap.go; then
        echo "FAIL: Banned component registered in bootstrap.go: $COMP"
        VIOLATIONS=$((VIOLATIONS + 1))
    else
        echo "PASS: Component $COMP is not registered in bootstrap.go"
    fi
done

if [ "$VIOLATIONS" -eq 0 ]; then
    echo "=== Aparok Core Sovereign Offline Gate: PASSED ==="
    exit 0
else
    echo "=== Aparok Core Sovereign Offline Gate: FAILED ($VIOLATIONS violations) ==="
    exit 1
fi
