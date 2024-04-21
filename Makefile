
build:
	docker build -t trevatk/mora:v0.0.1 -f ./docker/Dockerfile .

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run ./...