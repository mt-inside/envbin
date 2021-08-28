REPO:="docker.io/mtinside/envbin"
TAG:="latest"

lint:
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
	go test ./...

build-daemon: lint
	go build -ldflags "$(build/ldflags.sh)" ./cmd/daemon/...

build-client: lint
	go build -ldflags "$(build/ldflags.sh)" ./cmd/client/...

install: lint
	./deploy/git-hooks/install-local

run-server: lint
	go run -ldflags "$(build/ldflags.sh)" ./cmd/daemon/... serve

run-dump: lint
	go run -ldflags "$(build/ldflags.sh)" ./cmd/daemon/... dump

run-client: lint
	go run -ldflags "$(build/ldflags.sh)" ./cmd/client/...

package:
	docker buildx build --load -t {{REPO}}:{{TAG}} .

publish: package
	docker push {{REPO}}

run-docker: package
	docker run --rm --name envbin -p8080:8080 {{REPO}}:{{TAG}}
