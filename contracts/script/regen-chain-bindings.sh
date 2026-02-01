#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[fleming] exporting ABI..."
bash "$repo_root/contracts/script/export-abis.sh"

echo "[fleming] regenerating Go bindings..."
(
  cd "$repo_root"
  go generate ./pkg/protocol/chain
)

echo "[fleming] done"

