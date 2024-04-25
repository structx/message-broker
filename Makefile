
build image:
	docker build -t trevatk/mora:v0.0.1 .

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run ./...

utest:
	go test -v ./...