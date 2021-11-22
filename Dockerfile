ARG ARCH=
FROM ${ARCH}golang:1.17 as build

ARG VERSION=unknown

# libusb needed to build it
# won't statically link - tries to pull in libusb, which in turn relies on libudev functions, which libudev-dev doesn't help with...
# TODO: think! need a separate build profile for inside a container (go's +build stuff should go it, just define a tag) - things like usb don't work there anyway
# RUN apt update
# RUN apt install -y libusb-1.0-0-dev libudev-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Because we're building *in* a container for a container, there's no cross-OS-compilation; no need to specify GOOS
# Also because we take ARG ARCH and use buildx (invokes qemu), we always use the native compiler for any platform; never any need to specify GOARCH
# TODO: understand this better, then consider -tags netgo,osusergo etc to use go bits instead of c bits (which we statically link in here)
RUN go install -a -ldflags "-w $(build/ldflags.sh $VERSION)" ./cmd/daemon/...
# Dockerfile should take git info as an ARG, to be provided by
# * makefile on local
# * another step on GH actions
# * -> pass as an arg to ldflags, if not set, try to read locally (ie for native builds)


FROM gcr.io/distroless/base:latest AS run

ARG PORT=8080

COPY --from=build /go/bin/daemon /envbin
COPY *tpl /

EXPOSE $PORT
ENTRYPOINT ["/envbin"]
CMD ["serve"]
