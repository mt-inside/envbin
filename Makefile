.PHONY: build run
.DEFAULT_GOAL := run

Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
BuildTime := $(shell date +%Y-%m-%d_%H:%M:%S%z)
LDFLAGS := -X main.Version=$(Version) -X main.GitCommit=$(GitCommit) -X main.BuildTime=$(BuildTime)
LDFLAGS += -X github.com/mt-inside/envbin/pkg/data.Version=$(Version) -X github.com/mt-inside/envbin/pkg/data.GitCommit=$(GitCommit) -X github.com/mt-inside/envbin/pkg/data.BuildTime=$(BuildTime)

lint:
	golangci-lint run

build:
	go build -ldflags "$(LDFLAGS)" cmd/envbin.go

build-docker:
	GOOS=linux go build -o envbin-docker -ldflags "$(LDFLAGS)" cmd/envbin.go

run:
	go run -ldflags "$(LDFLAGS)" cmd/envbin.go :8088

image: build-docker
	docker build -t mtinside/envbin:latest .
image-run: image
	docker run -p8080:8080 mtinside/envbin:latest
image-push: image
	docker push mtinside/envbin
