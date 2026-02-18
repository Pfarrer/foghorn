FROM alpine:latest

RUN apk add --no-cache bash curl bc ca-certificates

COPY check.sh /check.sh
RUN chmod +x /check.sh

CMD ["/check.sh"]
