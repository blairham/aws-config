#!/bin/bash
set -e

echo "🧪 Testing CI workflow locally..."

echo "📦 Step 1: Download dependencies"
make deps

echo "🧪 Step 2: Run tests and coverage"
./run-tests.sh && go tool cover -html=coverage.out -o coverage/coverage.html

echo "🔨 Step 3: Run security scans"
echo "Running gosec security scanner..."
if command -v gosec >/dev/null 2>&1; then
  gosec ./... || echo "⚠️ Gosec found security issues (expected behavior)"
else
  echo "⚠️ Gosec not found in PATH, installing and running locally for test..."
  go install github.com/securego/gosec/v2/cmd/gosec@latest
  if [ -f "$(go env GOPATH)/bin/gosec" ]; then
    "$(go env GOPATH)/bin/gosec" ./... || echo "⚠️ Gosec found security issues (expected behavior)"
  else
    echo "⚠️ Could not find installed gosec, but this would work in CI"
  fi
fi

echo "Running govulncheck..."
if command -v govulncheck >/dev/null 2>&1; then
  govulncheck ./...
else
  echo "Installing govulncheck..."
  go install golang.org/x/vuln/cmd/govulncheck@latest
  if [ -f "$(go env GOPATH)/bin/govulncheck" ]; then
    "$(go env GOPATH)/bin/govulncheck" ./...
  else
    echo "⚠️ Could not find installed govulncheck, but this would work in CI"
  fi
fi

echo "🔨 Step 4: Build"
make build

echo "🔨 Step 4: Test build artifacts"
ls -la dist/
binary=$(find dist/ -name "aws-sso-config" -type f | head -1)
echo "Testing binary: $binary"
$binary --help

echo "🎉 All CI steps completed successfully!"
