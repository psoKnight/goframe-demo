#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

protoc -I . \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  calculator.proto
