ARG ARCH=
FROM ${ARCH}golang:1.16 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# including the .git dir
COPY . .
# Because we're building *in* a container for a container, there's no cross-OS-compilation; no need to specify GOOS
# Also because we take ARG ARCH and use buildx (invokes qemu), we always use the native compiler for any platform; never any need to specify GOARCH
RUN CGO_ENABLED=1 go install -a -tags netgo -ldflags "-w -extldflags '-static' $(build/ldflags.sh)" cmd/envbin.go


FROM gcr.io/distroless/static-debian10:latest AS run

ARG PORT=8080

COPY --from=build /go/bin/envbin /
COPY *tpl /

EXPOSE $PORT
ENTRYPOINT ["/envbin"]
CMD ["serve"]
