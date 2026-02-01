#!/usr/bin/env bash
set -euo pipefail

# This script lives at: pkg/protocol/chain/script/generate-bindings.sh
# Going up 4 directories lands at the repo root.
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../../" && pwd)"

abi_file="$repo_root/pkg/protocol/chain/abi/FlemingAnchor.abi.json"
out_file="$repo_root/pkg/protocol/chain/fleming_anchor_abigen.go"

# Keep Go tool caches within the repo to avoid permission issues in restricted
# environments (and to make CI/dev runs more predictable).
export GOPATH="${repo_root}/.cache/gopath"
export GOMODCACHE="${repo_root}/.cache/gopath/pkg/mod"
export GOCACHE="${repo_root}/.cache/gocache"

mkdir -p "$(dirname "$out_file")" "$GOMODCACHE" "$GOCACHE"

go run github.com/ethereum/go-ethereum/cmd/abigen@v1.16.8 \
  --abi "$abi_file" \
  --pkg chain \
  --type FlemingAnchor \
  --out "$out_file"

