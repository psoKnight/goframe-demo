#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
EXEFILE="generalgateway"
IMAGE_REPO="${IMAGE_REPO:-goframe-demo}"

usage() {
    echo "Usage:"
    echo "  $0 -b [test|online]"
    echo "  $0 -d VERSION [test|online]"
    echo "  $0 -b [test|online] delete"
    echo "  $0 -d VERSION [test|online] delete"
}

select_config() {
    case "$1" in
        online) echo "${SCRIPT_DIR}/config/config.yml" ;;
        test)   echo "${SCRIPT_DIR}/config/config_test.yml" ;;
        *) echo "Unknown environment: $1" >&2; exit 1 ;;
    esac
}

build_release() {
    local environment="$1"
    local config_file
    config_file="$(select_config "${environment}")"

    bash "${PROJECT_DIR}/pbdata/greeterpb/pb.sh"
    bash "${PROJECT_DIR}/pbdata/calculatorpb/pb.sh"
    rm -rf "${SCRIPT_DIR}/docker/release"
    mkdir -p "${SCRIPT_DIR}/docker/release/conf"

    (
        cd "${SCRIPT_DIR}"
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
            go build -trimpath -ldflags="-s -w" -o docker/release/${EXEFILE} .
    )

    cp "${SCRIPT_DIR}/start.sh" "${SCRIPT_DIR}/docker/release/start.sh"
    cp "${config_file}" "${SCRIPT_DIR}/docker/release/conf/config.yml"
    chmod +x "${SCRIPT_DIR}/docker/release/start.sh" "${SCRIPT_DIR}/docker/release/${EXEFILE}"
    echo "${EXEFILE} build success"
}

docker_login() {
    if [[ -n "${DOCKER_USERNAME:-}" ]]; then
        : "${DOCKER_PASSWORD:?DOCKER_PASSWORD is required when DOCKER_USERNAME is set}"
        printf '%s' "${DOCKER_PASSWORD}" | docker login "${DOCKER_REGISTRY:-}" \
            --username "${DOCKER_USERNAME}" --password-stdin
    fi
}

deploy_image() {
    local version="$1"
    local environment="$2"
    local image_name="${IMAGE_REPO}/${EXEFILE}"

    if [[ "${environment}" != "online" ]]; then
        image_name="${image_name}-${environment}"
    fi

    docker_login
    docker build -t "${image_name}:${version}" "${SCRIPT_DIR}/docker"
    docker push "${image_name}:${version}"
    echo "Image pushed: ${image_name}:${version}"
}

main() {
    local action="${1:-}"
    local delete_flag=""
    if [[ "${*: -1}" == "delete" ]]; then
        delete_flag="delete"
        set -- "${@:1:$(($# - 1))}"
    fi

    case "${action}" in
        -b)
            build_release "${2:-test}"
            ;;
        -d)
            [[ $# -ge 2 ]] || { usage; exit 1; }
            build_release "${3:-test}"
            deploy_image "${2}" "${3:-test}"
            ;;
        *)
            usage
            exit 1
            ;;
    esac

    if [[ "${delete_flag}" == "delete" ]]; then
        rm -rf "${SCRIPT_DIR}/docker/release"
        rm -f "${PROJECT_DIR}/pbdata/greeterpb/greeter.pb.go"
        rm -f "${PROJECT_DIR}/pbdata/greeterpb/greeter_grpc.pb.go"
        rm -f "${PROJECT_DIR}/pbdata/calculatorpb/calculator.pb.go"
        rm -f "${PROJECT_DIR}/pbdata/calculatorpb/calculator_grpc.pb.go"
        echo "${EXEFILE} release artifacts deleted"
    fi
}

main "$@"
