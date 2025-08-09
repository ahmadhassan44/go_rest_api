build:
	@go build -o bin/go_rest_api 

run: build
	@./bin/go_rest_api