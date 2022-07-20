export GOPATH=${HOME}/go
MAIN_FILES := $(shell find . ! -name "*_test.go" -name "*.go")

linter:
	${GOPATH}/bin/golangci-lint run .
run:
	go run $(MAIN_FILES)
