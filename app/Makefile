.PHONY: build test run

build:
	go build -o main .

test:
	go test -v ./...

run:
	go run main.go

docker-build:
	docker build -t myapp .

docker-run:
	docker run -p 8080:8080 myapp