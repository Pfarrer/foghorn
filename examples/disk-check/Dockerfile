FROM alpine:latest

RUN apk add --no-cache bash coreutils

COPY check.sh /check.sh
RUN chmod +x /check.sh

CMD ["/check.sh"]