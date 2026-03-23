#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "Building gocoding..."
go build -o ~/go/bin/gocoding ./cmd/gocoding/

echo "Installed to ~/go/bin/gocoding"
echo "Run 'gocoding' to start, or 'gocoding -p' for provider config"
