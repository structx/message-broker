
build:
	docker build -t trevatk/message-broker:latest -f ./build/Dockerfile .

deps:
	go mod tidy
	go mod vendor