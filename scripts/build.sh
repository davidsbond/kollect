#!/usr/bin/env bash

# This script iterates over each subdirectory in /cmd that contains a main.go file
# and builds the binary. Each binary is placed within /bin. Compilation details are
# linked in the binary via ldflags.

DIR=$(pwd)
BIN_DIR=${DIR}/bin
APP_VERSION=$(git describe --tags --always)
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
GOARM=$(go env GOARM)

rm -rf "${BIN_DIR}"
mkdir -p "${BIN_DIR}"

function compile_all() {
  compile_linux_amd64
  compile_linux_arm64
}

function compile_linux_arm64() {
  GOOS=linux
  GOARCH=arm64

  compile
}

function compile_linux_amd64() {
  GOOS=linux
  GOARCH=amd64

  compile
}

function compile() {
  APP_NAME="kollect"
  APP_DESCRIPTION="Monitor your Kubernetes clusters via your favourite event bus"
  COMPILED=$(date +%s)

  OUTPUT=${BIN_DIR}/${GOOS}/${GOARCH}/${APP_NAME}

  # If GOARM has been set, add an extra directory.
  if [ -n "$GOARM" ]; then
    OUTPUT=${BIN_DIR}/${GOOS}/${GOARCH}/v${GOARM}/${APP_NAME}
  fi

  GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} CGO_ENABLED=0 go build -ldflags \
    "-w -s "\
"-X 'github.com/davidsbond/kollect/internal/environment.Version=${APP_VERSION}'"\
"-X 'github.com/davidsbond/kollect/internal/environment.compiled=${COMPILED}'"\
"-X 'github.com/davidsbond/kollect/internal/environment.ApplicationName=${APP_NAME}'"\
"-X 'github.com/davidsbond/kollect/internal/environment.ApplicationDescription=${APP_DESCRIPTION}'" \
    -o "${OUTPUT}" .
}

compile_all
