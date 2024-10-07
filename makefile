build:
	@mkdir -p bin
	@go build -o bin/namecheap pkg/main.go

run: build
	@./bin/namecheap

lint:
	@golangci-lint run