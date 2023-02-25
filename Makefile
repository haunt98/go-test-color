.PHONY: all test-color lint try-4-real

all: test-color lint
	go mod tidy
	$(MAKE) test-color
	$(MAKE) lint
	$(MAKE) build
	$(MAKE) clean

test-color:
	go install github.com/haunt98/go-test-color@latest
	go-test-color -race -failfast .

lint:
	golangci-lint run .

try-4-real:
	go run . -race ./example/...

build:
	$(MAKE) clean
	go build -o go-test-color .

clean:
	rm -f go-test-color
