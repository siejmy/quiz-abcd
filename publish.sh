#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${DIR}"
set -e

BASE_TAG="jedrzejlewandowski/siejmy-quiz-abcd"

docker build -t "${BASE_TAG}" .
docker tag "${BASE_TAG}" "${BASE_TAG}:latest"
docker tag "${BASE_TAG}" "${BASE_TAG}:1.1.0"
docker push "${BASE_TAG}:latest"
docker push "${BASE_TAG}:1.1.0"
