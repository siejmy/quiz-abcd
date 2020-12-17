#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${DIR}"
set -e

BASE_TAG="jedrzejlewandowski/siejmy-quiz-abcd"
VERSION="2.5.0"

docker build -t "${BASE_TAG}" .
docker tag "${BASE_TAG}" "${BASE_TAG}:latest"
docker tag "${BASE_TAG}" "${BASE_TAG}:${VERSION}"
docker push "${BASE_TAG}:latest"
docker push "${BASE_TAG}:${VERSION}"
