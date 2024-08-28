#!/bin/bash

set -e

docker build --platform linux/arm/v8  -t benjamincaldwell/bento:latest-arm .
docker build --platform linux/amd64  -t benjamincaldwell/bento:latest-amd .
docker push benjamincaldwell/bento:latest-amd
docker push benjamincaldwell/bento:latest-arm
docker manifest create \
benjamincaldwell/bento:$1 \
--amend benjamincaldwell/bento:latest-amd \
--amend benjamincaldwell/bento:latest-arm
docker manifest push benjamincaldwell/bento:$1

