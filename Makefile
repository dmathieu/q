GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test: tidy
	go test -race -v ./...

tidy: goimports golint
	test -z "$$(goimports -l -d $(GO_FILES) | tee /dev/stderr)"
	test -z "$$(golint $(GO_FILES) | tee /dev/stderr)"
	test -z "$$(go vet $(GO_FILES) | tee /dev/stderr)"

golint:
	go get github.com/golang/lint/golint

goimports:
	go get golang.org/x/tools/cmd/goimports
