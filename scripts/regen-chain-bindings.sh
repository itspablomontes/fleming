#!/usr/bin/env bash
set -euo pipefail

# Regenerate Go bindings for the FlemingAnchor contract.
#
# Prerequisites:
#   - Foundry (forge) installed
#   - Go 1.21+ installed
#
# Usage:
#   bash scripts/regen-chain-bindings.sh

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
contracts_root="$repo_root/contracts"
abi_out_dir="$repo_root/pkg/protocol/chain/abi"
abi_out_file="$abi_out_dir/FlemingAnchor.abi.json"
go_out_file="$repo_root/apps/backend/internal/chain/fleming_anchor_abigen.go"

mkdir -p "$abi_out_dir"

echo "[fleming] Exporting ABI from Foundry (protocol source-of-truth)..."
bash "$repo_root/contracts/script/export-abis.sh"

echo "[fleming] Regenerating Go bindings..."
export GOPATH="${repo_root}/.cache/gopath"
export GOMODCACHE="${repo_root}/.cache/gopath/pkg/mod"
export GOCACHE="${repo_root}/.cache/gocache"
mkdir -p "$GOMODCACHE" "$GOCACHE"

go run github.com/ethereum/go-ethereum/cmd/abigen@v1.16.8 \
  --abi "$abi_out_file" \
  --pkg chain \
  --type FlemingAnchor \
  --out "$go_out_file"

echo "[fleming] Chain bindings regenerated successfully"
