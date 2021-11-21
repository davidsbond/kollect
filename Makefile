.DEFAULT_GOAL := go-build
TARGETPLATFORM := $(shell go env GOOS)/$(shell go env GOARCH)

go-modules:
	go mod tidy
	go mod vendor

go-install:
	go install \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		mvdan.cc/gofumpt/gofumports \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		github.com/bufbuild/buf/cmd/buf

go-lint:
	golangci-lint run --enable-all

go-format:
	grep -L -R "Code generated .* DO NOT EDIT" --exclude-dir=.git --exclude-dir=vendor --include="*.go" | \
	xargs -n 1 gofumports -w -local github.com/davidsbond/kollect

go-build:
	./scripts/build.sh

go-test:
	go test -race ./...

go-pack:
	upx `find ./bin -type f`

proto-generate:
	buf generate

proto-lint:
	buf lint

proto-breaking:
	buf breaking --against 'https://github.com/davidsbond/kollect.git#branch=master'

docker-build:
	./scripts/docker.sh

docker-compose:
	docker-compose down
	docker-compose build --build-arg TARGETPLATFORM=${TARGETPLATFORM}
	docker-compose up

has-changes:
	git add .
	git diff --staged --exit-code
