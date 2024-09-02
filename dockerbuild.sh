#!/bin/bash

set -e

docker build --platform linux/arm/v8  -t ${DOCKER_IMAGE}:latest-arm .
docker build --platform linux/amd64  -t ${DOCKER_IMAGE}:latest-amd .
docker push ${DOCKER_IMAGE}:latest-amd
docker push ${DOCKER_IMAGE}:latest-arm
docker manifest create \
${DOCKER_IMAGE}:$1 \
--amend ${DOCKER_IMAGE}:latest-amd \
--amend ${DOCKER_IMAGE}:latest-arm
docker manifest push ${DOCKER_IMAGE}:$1

curl -H "Authorization: Bearer $(gh auth token)" "https://api.github.com/repos/${K8S_CONFIG_REPO}/actions/workflows/update-image.yaml/dispatches" -d '{"ref": "main", "inputs": {"repository": "'"${DOCKER_IMAGE}"'", "tag": "'"${1}"'"}}'
