all: tidy test-color lint

tidy:
    go mod tidy

test-color:
    go install github.com/haunt98/go-test-color@latest
    go-test-color -race -failfast .

lint:
    golangci-lint run --fix ./...
    modernize -fix -test ./...

try-4-real:
    go run . -race ./example/...

build:
    go build -o go-test-color .

clean:
    rm -f go-test-color
