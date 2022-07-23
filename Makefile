.PHONY: all test-color lint try-4-real

all: test-color lint try-4-real

test-color:
	go install github.com/haunt98/go-test-color@latest
	go-test-color -race -failfast ./...

lint:
	golangci-lint run ./...

try-4-real:
	go run .
