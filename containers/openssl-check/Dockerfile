FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends openssl ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY check.sh /check.sh
RUN chmod +x /check.sh

CMD ["/check.sh"]
