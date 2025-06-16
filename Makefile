build:
	 @go build -o bin/deerDB cmd/main.go

run: build 
	@./bin/deerDB

test:
	@go test -v ./...
