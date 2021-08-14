REPO:="docker.io/mtinside/envbin"
TAG:="latest"

lint:
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
	go test ./...

build: lint
	go build -o envbin ./cmd/envbin/...

install: lint
	./deploy/git-hooks/install-local

run-server: lint
	go run -ldflags "$(build/ldflags.sh)" ./cmd/envbin/... serve

run-oneshot: lint
	go run -ldflags "$(build/ldflags.sh)" ./cmd/envbin/... oneshot

package:
	docker buildx build --load -t {{REPO}}:{{TAG}} .

publish: package
	docker push {{REPO}}

run-docker: package
	docker run --rm --name envbin -p8080:8080 {{REPO}}:{{TAG}}
