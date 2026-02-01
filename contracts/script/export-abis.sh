#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
contracts_root="$repo_root/contracts"

abi_out_dir="$repo_root/pkg/protocol/chain/abi"
abi_out_file="$abi_out_dir/FlemingAnchor.abi.json"

mkdir -p "$abi_out_dir"

cd "$contracts_root"

# Source of truth: Foundry compilation output for the contract ABI.
# We commit this JSON (not `out/` artifacts) so downstream tooling (abigen, docs, CI)
# has a stable, reviewable input without requiring Foundry at build/runtime.
forge inspect FlemingAnchor abi --json > "$abi_out_file"

