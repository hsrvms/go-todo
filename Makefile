run: build
	@./bin/api
build: 
	@go build -o bin/api
test:
	@go test -v -cover ./...
docs:
	godoc -http :8081
