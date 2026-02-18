#!/bin/sh

set -e

START_TIME=$(date +%s%3N)

URL="${CHECK_URL:-https://example.com}"
EXPECTED_STATUS="${EXPECTED_STATUS:-200}"
TIMEOUT_RAW="${REQUEST_TIMEOUT:-10}"
VERIFY_SSL_RAW="$(echo "${VERIFY_SSL:-true}" | tr '[:upper:]' '[:lower:]')"

case "$TIMEOUT_RAW" in
    *s)
        TIMEOUT="${TIMEOUT_RAW%s}"
        ;;
    *m)
        TIMEOUT="$(( ${TIMEOUT_RAW%m} * 60 ))"
        ;;
    *h)
        TIMEOUT="$(( ${TIMEOUT_RAW%h} * 3600 ))"
        ;;
    *)
        TIMEOUT="$TIMEOUT_RAW"
        ;;
esac

case "$TIMEOUT" in
    ''|*[!0-9]*)
        TIMEOUT="10"
        ;;
esac

SSL_FLAG=""
case "$VERIFY_SSL_RAW" in
    true|1|yes)
        ;;
    false|0|no)
        SSL_FLAG="-k"
        ;;
    *)
        VERIFY_SSL_RAW="true"
        ;;
esac

HTTP_CODE=$(curl -sS -o /dev/null -w "%{http_code}" --max-time "$TIMEOUT" $SSL_FLAG "$URL" 2>/dev/null || true)
case "$HTTP_CODE" in
    ''|*[!0-9]*)
        HTTP_CODE="000"
        ;;
esac

END_TIME=$(date +%s%3N)
DURATION=$((END_TIME - START_TIME))

if [ "$HTTP_CODE" = "$EXPECTED_STATUS" ]; then
    STATUS="pass"
    MESSAGE="HTTP check passed with status $HTTP_CODE"
else
    STATUS="fail"
    MESSAGE="HTTP check failed: expected $EXPECTED_STATUS, got $HTTP_CODE"
fi

cat <<RESULT
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "url": "$URL",
    "http_code": "$HTTP_CODE",
    "expected_status": "$EXPECTED_STATUS"
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
RESULT
