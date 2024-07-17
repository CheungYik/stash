build:
	@go build -o bin/stash

run: build
	@./bin/stash

test:
	@go test -v ./...
