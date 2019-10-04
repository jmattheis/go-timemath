download-tools:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

lint:
	go vet ./...
	golint -set_exit_status $(shell go list ./...)
	goimports -l $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...

install:
	go mod download
