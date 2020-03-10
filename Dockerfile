FROM gcr.io/distroless/base-debian10:latest

ARG PORT=8080

WORKDIR /app
COPY envbin-docker envbin
COPY *tpl ./

EXPOSE $PORT
ENTRYPOINT ["/app/envbin"]
CMD ["--port", "8080"]
