.PHONY: lint run image image-push docker-run
.DEFAULT_GOAL := run-oneshot


lint:
	#go fmt ./...
	#go vet ./...
	#golangci-lint run ./...

build: lint
	go build -o envbin ./cmd/envbin/...

install: lint
	./deploy/git-hooks/install-local

run-server: lint
	go run -ldflags "$(shell build/ldflags.sh)" ./cmd/envbin/... serve

run-oneshot: lint
	go run -ldflags "$(shell build/ldflags.sh)" ./cmd/envbin/... oneshot

image:
	docker buildx build --load -t mtinside/envbin:latest .
image-push: image
	docker push mtinside/envbin
docker-run: image
	docker run --rm --name envbin -p8080:8080 mtinside/envbin:latest
