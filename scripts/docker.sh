#!/usr/bin/env bash

REGISTRY=ghcr.io
IMAGE=${REGISTRY}/davidsbond/kollect
TAG=$(git describe --tags --always)

echo "${GITHUB_TOKEN}" | docker login ${REGISTRY} -u davidsbond --password-stdin

docker buildx create --use
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --push \
  --label "org.opencontainers.image.created=$(date --rfc-3339=seconds)" \
  --label "org.opencontainers.image.authors=David Bond <davidsbond93@gmail.com>" \
  --label "org.opencontainers.image.url=https://github.com/davidsbond/kollect" \
  --label "org.opencontainers.image.documentation=https://github.com/davidsbond/kollect/blob/master/README.md" \
  --label "org.opencontainers.image.source=https://github.com/davidsbond/kollect" \
  --label "org.opencontainers.image.version=${TAG}" \
  --label "org.opencontainers.image.revision=${TAG}" \
  --label "org.opencontainers.image.vendor=David Bond" \
  --label "org.opencontainers.image.licenses=Apache-2.0" \
  --label "org.opencontainers.image.title=Kollect" \
  --label "org.opencontainers.image.description=Monitor your Kubernetes clusters via your favourite event bus" \
  --label "org.opencontainers.image.base.name=gcr.io/distroless/static" \
  -t ${IMAGE}:latest \
  -t ${IMAGE}:"${TAG}" .
