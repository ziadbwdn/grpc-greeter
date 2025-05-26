#!/bin/bash

set -e

PROJECT_DIR=$(pwd)
PROTO_DIR="$PROJECT_DIR/pkg/proto"
OUT_DIR="$PROTO_DIR/generated"

PROTOC_GEN_GO=$(which protoc-gen-go || true)
PROTOC_GEN_GO_GRPC=$(which protoc-gen-go-grpc || true)

log() {
    echo -e "[INFO] $1"
}

err() {
    echo -e "[ERROR] $1" >&2
    exit 1
}

log "Working in project: $PROJECT_DIR"

# Check protoc
if ! command -v protoc &>/dev/null; then
    err "protoc not found. Please install Protocol Buffers compiler."
fi

log "Found protoc version: $(protoc --version)"

# Check plugins
if [ -z "$PROTOC_GEN_GO" ]; then
    err "protoc-gen-go is not installed.\nInstall it with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
fi

if [ -z "$PROTOC_GEN_GO_GRPC" ]; then
    err "protoc-gen-go-grpc is not installed.\nInstall it with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
fi

log "Found protoc-gen-go at: $PROTOC_GEN_GO"
log "Found protoc-gen-go-grpc at: $PROTOC_GEN_GO_GRPC"

mkdir -p "$OUT_DIR"

log "Generating gRPC code..."
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$OUT_DIR" \
    --go-grpc_out="$OUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/greet.proto"

log "Protobuf generation completed."
