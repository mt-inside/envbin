.PHONY: build run
.DEFAULT_GOAL := run

FLAGS := -ldflags "-X data.version=0.0.1"

lint:
	golangci-lint run

build:
	go build $(FLAGS) cmd/envbin.go

build-docker:
	GOOS=linux go build -o envbin-docker $(FLAGS) cmd/envbin.go

run:
	go run $(FLAGS) cmd/envbin.go

image: build-docker
	docker build -t mtinside/envbin:latest .
image-run: image
	docker run -p8080:8080 mtinside/envbin:latest
image-push: image
	docker push mtinside/envbin
