FROM golang:1.15 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# including the .git dir
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags "-w -extldflags '-static' $(build/ldflags.sh)" cmd/envbin.go


FROM gcr.io/distroless/static-debian10:latest AS run

ARG PORT=8080

COPY --from=build /go/bin/envbin /
COPY *tpl /

EXPOSE $PORT
ENTRYPOINT ["/envbin"]
