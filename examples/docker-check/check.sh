#!/bin/bash

set -e

START_TIME=$(date +%s%3N)

CHECK_NAME="${FOGHORN_CHECK_NAME:-unknown}"
ENDPOINT="${FOGHORN_ENDPOINT:-http://localhost:8080/health}"
TIMEOUT="${FOGHORN_TIMEOUT:-10s}"

echo "Running check: $CHECK_NAME" >&2
echo "Endpoint: $ENDPOINT" >&2

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "$ENDPOINT" 2>&1 || echo "000")

END_TIME=$(date +%s%3N)
DURATION=$((END_TIME - START_TIME))

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "201" ]; then
    STATUS="pass"
    MESSAGE="Health check passed with HTTP $HTTP_CODE"
else
    STATUS="fail"
    MESSAGE="Health check failed with HTTP $HTTP_CODE"
fi

cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "http_code": $HTTP_CODE,
    "endpoint": "$ENDPOINT"
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
EOF
