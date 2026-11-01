#!/bin/sh
set -eu
cd "$(dirname "$0")/../api"
task_generated=$(mktemp)
trap 'rm -f "$task_generated"' EXIT HUP INT TERM
go tool oapi-codegen -config oapi-codegen.yaml -o "$task_generated" openapi.yaml
if ! cmp -s ../internal/api/api.gen.go "$task_generated"; then
    printf '%s\n' 'Generated API is stale. Run: go generate ./api' >&2
    exit 1
fi
