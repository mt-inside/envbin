FROM gcr.io/distroless/base-debian10:latest

ARG PORT=8080

COPY envbin2-docker /app/envbin2

EXPOSE $PORT
ENTRYPOINT ["/app/envbin2"]
CMD ["--port", "8080"]
