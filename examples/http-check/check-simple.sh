#!/bin/sh

START_TIME=$(date +%s%3N)
URL="${URL:-https://example.com}"
EXPECTED_STATUS="${EXPECTED_STATUS:-200}"

echo "Testing HTTP to $URL" >&2

# Simple curl test
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "$URL" 2>&1 || echo "000")

END_TIME=$(date +%s%3N)
DURATION=$((END_TIME - START_TIME))

if [ "$HTTP_CODE" = "$EXPECTED_STATUS" ]; then
    STATUS="pass"
    MESSAGE="HTTP check passed with status $HTTP_CODE"
else
    STATUS="fail"
    MESSAGE="HTTP check failed: expected $EXPECTED_STATUS, got $HTTP_CODE"
fi

cat <<EOF
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
EOF
