.PHONY: lint run image image-push docker-run
.DEFAULT_GOAL := run


lint:
	golangci-lint run

run-server:
	go run -ldflags "$(shell build/ldflags.sh)" cmd/envbin.go serve

run-oneshot:
	go run -ldflags "$(shell build/ldflags.sh)" cmd/envbin.go oneshot

image:
	docker build -t mtinside/envbin:latest .
image-push: image
	docker push mtinside/envbin
docker-run: image
	docker run --rm --name envbin -p8080:8080 mtinside/envbin:latest
