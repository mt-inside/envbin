.PHONY: build run
.DEFAULT_GOAL := run


lint:
	golangci-lint run

run:
	go run -ldflags "$(shell build/ldflags.sh)" cmd/envbin.go

image:
	docker build -t mtinside/envbin:latest .
image-push: image
	docker push mtinside/envbin
docker-run: image
	docker run --rm --name envbin -p8080:8080 mtinside/envbin:latest
