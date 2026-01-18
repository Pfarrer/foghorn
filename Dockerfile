FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

ARG BINARY=foghorn
COPY ${BINARY} .

USER 65532:65532

ENTRYPOINT ["/app/foghorn"]
