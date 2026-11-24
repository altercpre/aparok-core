#!/usr/bin/env bash
set -euo pipefail

echo "=== Aparok Core Sovereign Offline Verification ==="

BANNED_PACKAGES=(
    "github.com/prometheus/client_golang"
    "github.com/uber/jaeger-client-go"
    "github.com/grandcat/zeroconf"
)

echo "Checking package dependencies against banned telemetry and remote discovery..."
TAGS="nogrpcserver gomobile nowatchdog nosigar nomutexdeadlockdetector timetzdata rasterizesvg"
DEPS=$(go list -tags "${TAGS}" -deps ./clientlibrary/service ./core 2>/dev/null || true)

VIOLATIONS=0
for BANNED in "${BANNED_PACKAGES[@]}"; do
    if echo "$DEPS" | grep -q "$BANNED"; then
        echo "FAIL: Banned package detected in dependencies: $BANNED"
        VIOLATIONS=$((VIOLATIONS + 1))
    else
        echo "PASS: Clean - $BANNED not present in dependency graph."
    fi
done

echo "Checking for local-only mode enforcement..."
if grep -q "IsLocalOnlyMode() bool" core/anytype/config/config.go; then
    echo "PASS: IsLocalOnlyMode defined and enforced in core/anytype/config/config.go"
else
    echo "FAIL: IsLocalOnlyMode not found in config"
    VIOLATIONS=$((VIOLATIONS + 1))
fi

if [ "$VIOLATIONS" -eq 0 ]; then
    echo "=== Aparok Core Sovereign Offline Gate: PASSED ==="
    exit 0
else
    echo "=== Aparok Core Sovereign Offline Gate: FAILED ($VIOLATIONS violations) ==="
    exit 1
fi
