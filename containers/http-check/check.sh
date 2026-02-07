#!/bin/sh

set -e

START_TIME=$(date +%s%3N)

URL="${FOGHORN_ENDPOINT:-${URL:-https://example.com}}"
EXPECTED_STATUS="${EXPECTED_STATUS:-200}"
TIMEOUT="${FOGHORN_TIMEOUT:-10s}"

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time "$TIMEOUT" "$URL" 2>&1 || echo "000")

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
    "http_code": $HTTP_CODE,
    "expected_status": "$EXPECTED_STATUS"
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
RESULT
