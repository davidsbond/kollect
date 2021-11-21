//go:build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "mvdan.cc/gofumpt/gofumports"
)
