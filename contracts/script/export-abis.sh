#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
contracts_root="$repo_root/contracts"

abi_out_dir="$repo_root/pkg/protocol/chain/abi"
abi_out_file="$abi_out_dir/FlemingAnchor.abi.json"

mkdir -p "$abi_out_dir"

cd "$contracts_root"

# Soldeer installs dependencies into `contracts/dependencies/`, which is intentionally gitignored.
# In CI (and fresh checkouts), that folder won't exist until we install deps.
deps_root="$contracts_root/dependencies"
oz_marker="$deps_root/@openzeppelin-contracts-5.2.0/access/Ownable.sol"
forge_std_marker="$deps_root/forge-std-1.9.7/src/Test.sol"

if [[ ! -f "$oz_marker" || ! -f "$forge_std_marker" ]]; then
  echo "[fleming] Installing Soldeer dependencies..." >&2
  forge soldeer install
fi

# Source of truth: Foundry compilation output for the contract ABI.
# We commit this JSON (not `out/` artifacts) so downstream tooling (abigen, docs, CI)
# has a stable, reviewable input without requiring Foundry at build/runtime.
forge inspect FlemingAnchor abi --json > "$abi_out_file"

