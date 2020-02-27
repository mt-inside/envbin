FROM gcr.io/distroless/base-debian10:latest

ARG PORT=8080

WORKDIR /app
COPY envbin2-docker envbin2
COPY *tpl ./

EXPOSE $PORT
ENTRYPOINT ["/app/envbin2"]
CMD ["--port", "8080"]
