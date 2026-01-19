#!/bin/bash

set -e

START_TIME=$(date +%s%3N)

# Configuration
URL="${URL}"
EXPECTED_STATUS="${EXPECTED_STATUS:-200}"
METHOD="${METHOD:-GET}"
HEADERS="${HEADERS:-}"
REQUEST_BODY="${REQUEST_BODY:-}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-30}"
FOLLOW_REDIRECTS="${FOLLOW_REDIRECTS:-true}"
VERIFY_SSL="${VERIFY_SSL:-true}"
WARNING_THRESHOLD_MS="${WARNING_THRESHOLD_MS:-1000}"
CRITICAL_THRESHOLD_MS="${CRITICAL_THRESHOLD_MS:-5000}"
CONTENT_REGEX="${CONTENT_REGEX:-}"
FOGHORN_TIMEOUT="${FOGHORN_TIMEOUT:-10}"

# Validate URL
if [ -z "$URL" ]; then
    STATUS="fail"
    MESSAGE="URL environment variable is required"
    cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {},
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $(($(date +%s%3N) - START_TIME))
}
EOF
    exit 1
fi

# Build curl command
CURL_CMD="curl -s -o /dev/null -w '%{http_code}\n%{size_download}\n%{time_total}\n' -X $METHOD"

# Set timeout
CURL_CMD="$CURL_CMD --max-time $TIMEOUT_SECONDS"

# Handle redirects
if [ "$FOLLOW_REDIRECTS" = "false" ]; then
    CURL_CMD="$CURL_CMD --location-trusted"
else
    CURL_CMD="$CURL_CMD -L"
fi

# Handle SSL verification
if [ "$VERIFY_SSL" = "false" ]; then
    CURL_CMD="$CURL_CMD -k"
fi

# Add headers
if [ -n "$HEADERS" ]; then
    # Parse JSON headers and add to curl command
    while IFS= read -r line; do
        # Skip opening/closing braces and commas
        if [[ "$line" =~ ^[[:space:]]*\} ]] || [[ "$line" =~ ^[[:space:]]*\{ ]] || [[ "$line" =~ ^[[:space:]]*$ ]]; then
            continue
        fi
        # Extract key and value
        key=$(echo "$line" | sed 's/^[[:space:]]*//' | cut -d'"' -f2)
        value=$(echo "$line" | sed 's/^[[:space:]]*//' | cut -d'"' -f4)
        if [ -n "$key" ]; then
            CURL_CMD="$CURL_CMD -H '$key: $value'"
        fi
    done <<< "$HEADERS"
fi

# Add request body for POST/PUT
if [ -n "$REQUEST_BODY" ]; then
    CURL_CMD="$CURL_CMD -d '$REQUEST_BODY'"
fi

# Add URL
CURL_CMD="$CURL_CMD '$URL'"

# Execute curl command
if ! OUTPUT=$(eval $CURL_CMD 2>&1); then
    END_TIME=$(date +%s%3N)
    DURATION=$((END_TIME - START_TIME))

    STATUS="fail"
    MESSAGE="Failed to connect to $URL: $OUTPUT"
    
    cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "url": "$URL",
    "method": "$METHOD",
    "error": "$OUTPUT"
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
EOF
    exit 1
fi

# Parse output
STATUS_CODE=$(echo "$OUTPUT" | head -1)
RESPONSE_SIZE=$(echo "$OUTPUT" | sed -n '2p')
RESPONSE_TIME=$(echo "$OUTPUT" | sed -n '3p')

# Convert response time to milliseconds
RESPONSE_TIME_MS=$(echo "$RESPONSE_TIME * 1000" | bc | awk '{printf "%.0f", $1}')

# Validate status code
STATUS_MATCH=0
if echo "$EXPECTED_STATUS" | grep -q '-'; then
    # Range of status codes (e.g., 200-299)
    MIN_CODE=$(echo "$EXPECTED_STATUS" | cut -d'-' -f1)
    MAX_CODE=$(echo "$EXPECTED_STATUS" | cut -d'-' -f2)
    if [ "$STATUS_CODE" -ge "$MIN_CODE" ] && [ "$STATUS_CODE" -le "$MAX_CODE" ]; then
        STATUS_MATCH=1
    fi
elif echo "$EXPECTED_STATUS" | grep -q ','; then
    # List of status codes (e.g., 200,201,204)
    IFS=',' read -ra CODES <<< "$EXPECTED_STATUS"
    for CODE in "${CODES[@]}"; do
        if [ "$STATUS_CODE" = "$CODE" ]; then
            STATUS_MATCH=1
            break
        fi
    done
else
    # Single status code
    if [ "$STATUS_CODE" = "$EXPECTED_STATUS" ]; then
        STATUS_MATCH=1
    fi
fi

# Check content regex if specified
CONTENT_MATCH=true
if [ -n "$CONTENT_REGEX" ]; then
    CONTENT=$(curl -s -X $METHOD --max-time $TIMEOUT_SECONDS "$URL")
    if [ "$FOLLOW_REDIRECTS" = "false" ]; then
        CONTENT=$(curl -s -X $METHOD --max-time $TIMEOUT_SECONDS --location-trusted "$URL")
    else
        CONTENT=$(curl -s -X $METHOD --max-time $TIMEOUT_SECONDS -L "$URL")
    fi

    if ! echo "$CONTENT" | grep -qE "$CONTENT_REGEX"; then
        CONTENT_MATCH=false
    fi
fi

END_TIME=$(date +%s%3N)
DURATION=$((END_TIME - START_TIME))

# Determine status
STATUS="pass"
if [ "$STATUS_MATCH" -eq 0 ]; then
    STATUS="fail"
elif [ "$CONTENT_MATCH" = "false" ]; then
    STATUS="fail"
elif [ "$RESPONSE_TIME_MS" -ge "$CRITICAL_THRESHOLD_MS" ]; then
    STATUS="fail"
elif [ "$RESPONSE_TIME_MS" -ge "$WARNING_THRESHOLD_MS" ]; then
    STATUS="warn"
fi

# Build message
STATUS_TEXT=$(curl -s -o /dev/null -w '%{http_code} %{remote_ip}' -X $METHOD --max-time $TIMEOUT_SECONDS "$URL" 2>/dev/null || echo "$STATUS_CODE")
if [ "$STATUS" = "pass" ]; then
    MESSAGE="HTTP $METHOD $URL: $STATUS_CODE ($RESPONSE_TIME_MS ms)"
elif [ "$STATUS" = "fail" ]; then
    if [ "$STATUS_MATCH" -eq 0 ]; then
        MESSAGE="HTTP $METHOD $URL: expected status $EXPECTED_STATUS, got $STATUS_CODE ($RESPONSE_TIME_MS ms)"
    elif [ "$CONTENT_MATCH" = "false" ]; then
        MESSAGE="HTTP $METHOD $URL: content regex '$CONTENT_REGEX' did not match ($RESPONSE_TIME_MS ms)"
    else
        MESSAGE="HTTP $METHOD $URL: response time $RESPONSE_TIME_MS ms exceeded critical threshold $CRITICAL_THRESHOLD_MS ms"
    fi
else
    MESSAGE="HTTP $METHOD $URL: $STATUS_CODE (response time $RESPONSE_TIME_MS ms above warning threshold $WARNING_THRESHOLD_MS ms)"
fi

# Output JSON
cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "url": "$URL",
    "method": "$METHOD",
    "status_code": $STATUS_CODE,
    "status_text": "$STATUS_CODE",
    "response_time_ms": $RESPONSE_TIME_MS,
    "response_size_bytes": ${RESPONSE_SIZE:-0},
    "content_match": $CONTENT_MATCH
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
EOF

# Exit with appropriate code
if [ "$STATUS" = "pass" ]; then
    exit 0
elif [ "$STATUS" = "unknown" ]; then
    exit 1
else
    exit 1
fi
