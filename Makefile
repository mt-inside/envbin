.PHONY: build run
.DEFAULT_GOAL := run

FLAGS := -ldflags "-X main.version=0.0.1"

build:
	go build $(FLAGS) cmd/envbin2.go

build-docker:
	GOOS=linux go build -o envbin2-docker $(FLAGS) cmd/envbin2.go

run:
	go run $(FLAGS) cmd/envbin2.go

image: build-docker
	docker build -t mtinside/envbin2:latest .
