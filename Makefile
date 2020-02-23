.PHONY: build run
.DEFAULT_GOAL := run

FLAGS := -ldflags "-X main.version=0.0.1"

build:
	go build $(FLAGS) ...

run:
	go run $(FLAGS) cmd/envbin2.go

