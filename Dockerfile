FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

ARG BINARY=foghorn-daemon
COPY --chown=65532:65532 ${BINARY} ./foghorn
RUN chmod +x ./foghorn

USER 65532:65532

ENTRYPOINT ["/app/foghorn"]
