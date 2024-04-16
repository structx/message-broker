
build:
	docker build -t trevatk/message-broker:latest -f ./docker/Dockerfile .

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run ./...