.PHONY: build test run

build:
	go build -o main .

test:
	go test -v ./...

run:
	go run main.go

docker-build:
	docker build -t app .

docker-run:
	docker run -p 8080:8080 app

update-docs:
	swag init