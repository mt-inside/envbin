set dotenv-load

DH_USER := "mtinside"
REPO:="docker.io/" + DH_USER + "/envbin"
TAG:=`git describe --tags --abbrev`
TAGD:=`git describe --tags --abbrev --dirty`
CGR_ARCHS := "amd64" # "amd64,aarch64,armv7"

default:
	@just --list

install-tools:
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

lint:
	go fmt ./...
	go vet ./...
	staticcheck -tags native ./...
	golangci-lint run --build-tags native ./...
	go test ./...

build-daemon: lint
	go build -tags native -ldflags "$(build/ldflags.sh)" ./cmd/daemon

build-client: lint
	go build -ldflags "$(build/ldflags.sh)" ./cmd/client

install: lint
	./deploy/git-hooks/install-local

run-server: lint
	go run -tags native -ldflags "$(build/ldflags.sh)" ./cmd/daemon serve
run-server-root: lint build-daemon
	sudo ./daemon serve

run-dump: lint
	go run -tags native -ldflags "$(build/ldflags.sh)" ./cmd/daemon dump
run-dump-root: lint build-daemon
	sudo ./daemon dump

run-client: lint
	go run -ldflags "$(build/ldflags.sh)" ./cmd/client

melange:
	# keypair to verify the package between melange and apko. apko will very quietly refuse to find our apk if these args aren't present
	docker run --rm -v "${PWD}":/work cgr.dev/chainguard/melange keygen
	docker run --privileged --rm -v "${PWD}":/work cgr.dev/chainguard/melange build --arch {{CGR_ARCHS}} --signing-key melange.rsa melange.yaml

package-cgr: melange
	docker run --rm -v "${PWD}":/work cgr.dev/chainguard/apko build -k melange.rsa.pub --build-arch {{CGR_ARCHS}} apko.yaml {{REPO}}:{{TAG}} envbin.tar
	docker load < envbin.tar
publish-cgr: melange
	docker run --rm -v "${PWD}":/work --entrypoint sh cgr.dev/chainguard/apko -c \
	'echo "'${DH_TOKEN}'" | apko login docker.io -u {{DH_USER}} --password-stdin && \
	apko publish apko.yaml {{REPO}}:{{TAG}} -k melange.rsa.pub --arch {{CGR_ARCHS}}'

package:
	docker buildx build --load -t {{REPO}}:{{TAG}} .

publish: package
	docker push {{REPO}}

run-docker: package
	docker run --rm --name envbin -p8080:8080 {{REPO}}:{{TAG}}
