#!/bin/bash
set -e

IMAGE_TAG="${IMAGE_TAG:-latest}"
RPC_IMAGE="wallet-rpc:${IMAGE_TAG}"
API_IMAGE="wallet-api:${IMAGE_TAG}"

build_rpc() {
    echo "Building ${RPC_IMAGE} ..."
    docker build -f Dockerfile.rpc -t "${RPC_IMAGE}" .
    echo "Done: ${RPC_IMAGE}"
}

build_api() {
    echo "Building ${API_IMAGE} ..."
    docker build -f Dockerfile.api -t "${API_IMAGE}" .
    echo "Done: ${API_IMAGE}"
}

case "${1}" in
    rpc)
        build_rpc
        ;;
    api)
        build_api
        ;;
    "")
        build_rpc
        build_api
        echo "All images built successfully."
        ;;
    *)
        echo "Usage: $0 [rpc|api]"
        echo "  (no args)  Build all images"
        echo "  rpc        Build wallet-rpc only"
        echo "  api        Build wallet-api only"
        echo ""
        echo "Environment variables:"
        echo "  IMAGE_TAG  Image tag (default: latest)"
        exit 1
        ;;
esac
