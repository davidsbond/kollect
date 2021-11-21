.DEFAULT_GOAL := go-build

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
	CGO_ENABLED=0 go build

go-test:
	go test -race ./...

proto-generate:
	buf generate

proto-lint:
	buf lint

proto-breaking:
	buf breaking --against 'https://github.com/davidsbond/kollect.git#branch=master'

docker-compose:
	docker-compose down
	docker-compose build
	docker-compose up

has-changes:
	git add .
	git diff --staged --exit-code
