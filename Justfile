set dotenv-load

default:
	@just --list --unsorted --color=always

REPO := "envbin"
CMD := "daemon"
DH_USER := "mtinside"
GH_USER := "mt-inside"
DH_REPO:="docker.io/" + DH_USER + "/envbin"
GH_REPO := "ghcr.io/" + GH_USER + "/print-cert"
TAG:=`git describe --tags --always --abbrev`
TAGD:=`git describe --tags --always --abbrev --dirty --broken`
CGR_ARCHS := "aarch64,amd64" # "x86,armv7"
LD_COMMON := "-ldflags \"-X 'github.com/mt-inside/envbin/pkg/data.Version=" + TAGD + "'\""
LD_STATIC := "-ldflags \"-X 'github.com/mt-inside/envbin/pkg/data.Version=" + TAGD + "' -w -linkmode external -extldflags '-static'\""
MELANGE := "melange"
APKO    := "apko"

tools-install:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/exp/cmd/...@latest
	go install github.com/kisielk/godepgraph@latest
	go install golang.org/x/tools/cmd/stringer@latest

lint:
	gofmt -s -w .
	goimports -local github.com/mt-inside/envbin -w .
	go vet ./...
	staticcheck -tags native ./...
	golangci-lint run --build-tags native ./...

test: lint
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

render-mod-graph:
	go mod graph | modgraphviz | dot -Tpng -o mod_graph.png

render-pkg-graph:
	godepgraph -s -onlyprefixes github.com/mt-inside ./cmd/daemon | dot -Tpng -o pkg_graph.png

build-daemon: lint
	CGO_ENABLED=0 go build -tags native {{LD_COMMON}} ./cmd/daemon

build-client: lint
	CGO_ENABLED=0 go build {{LD_COMMON}} ./cmd/client

install: lint
	CGO_ENABLED=0 ./deploy/git-hooks/install-local

run-server: lint
	CGO_ENABLED=0 go run -tags native {{LD_COMMON}} ./cmd/daemon serve
run-server-root: lint build-daemon
	sudo ./daemon serve

run-dump: lint
	CGO_ENABLED=0 go run -tags native {{LD_COMMON}} ./cmd/daemon dump
run-dump-root: lint build-daemon
	sudo ./daemon dump

run-client: lint
	CGO_ENABLED=0 go run {{LD_COMMON}} ./cmd/client

package: test
	# if there's >1 package in this directory, apko seems to pick the _oldest_ without fail
	rm -rf ./packages/
	{{MELANGE}} bump melange.yaml {{TAGD}}
	{{MELANGE}} keygen
	{{MELANGE}} build --arch {{CGR_ARCHS}} --signing-key melange.rsa melange.yaml

image-local:
	{{APKO}} build --keyring-append melange.rsa.pub --arch {{CGR_ARCHS}} apko.yaml {{GH_REPO}}:{{TAG}} {{CMD}}.tar
	docker load < {{CMD}}.tar

image-publish:
	{{APKO}} login docker.io -u {{DH_USER}} --password "${DH_TOKEN}"
	{{APKO}} login ghcr.io   -u {{GH_USER}} --password "${GH_TOKEN}"
	{{APKO}} publish --keyring-append melange.rsa.pub --arch {{CGR_ARCHS}} apko.yaml {{GH_REPO}}:{{TAG}} {{DH_REPO}}:{{TAG}}
