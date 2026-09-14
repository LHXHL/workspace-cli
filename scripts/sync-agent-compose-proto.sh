#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source_file="$repo_root/products/agentcompose/proto/SOURCE"

source "$source_file"
core_repo="$repo_root/$repository"
snapshot_dir="$repo_root/products/agentcompose/proto"

if ! git -C "$core_repo" cat-file -e "$commit^{commit}"; then
  echo "agent-compose source commit is unavailable: $commit" >&2
  exit 1
fi

mkdir -p "$snapshot_dir/agentcompose/v2" "$snapshot_dir/health/v1"
git -C "$core_repo" show "$commit:proto/agentcompose/v2/agentcompose.proto" > "$snapshot_dir/agentcompose/v2/agentcompose.proto"
git -C "$core_repo" show "$commit:proto/health/v1/health.proto" > "$snapshot_dir/health/v1/health.proto"

actual_agentcompose="$(shasum -a 256 "$snapshot_dir/agentcompose/v2/agentcompose.proto" | awk '{print $1}')"
actual_health="$(shasum -a 256 "$snapshot_dir/health/v1/health.proto" | awk '{print $1}')"

if [[ "$actual_agentcompose" != "$agentcompose_sha256" ]]; then
  echo "agentcompose proto checksum mismatch: got $actual_agentcompose" >&2
  exit 1
fi
if [[ "$actual_health" != "$health_sha256" ]]; then
  echo "health proto checksum mismatch: got $actual_health" >&2
  exit 1
fi

echo "synced Agent Compose proto snapshot at $commit"
