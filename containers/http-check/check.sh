#!/bin/sh

set -e

now_ms() {
    raw="$(date +%s%3N 2>/dev/null || true)"
    case "$raw" in
        ''|*[!0-9]*)
            echo $(( $(date +%s) * 1000 ))
            ;;
        *)
            echo "$raw"
            ;;
    esac
}

json_escape() {
    printf '%s' "$1" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g' -e 's/\t/\\t/g' -e ':a;N;$!ba;s/\n/\\n/g'
}

json_output() {
    status="$1"
    message="$2"
    duration_ms="$3"
    data="$4"

    escaped_msg="$(json_escape "$message")"

    cat <<RESULT
{
  "status": "$status",
  "message": "$escaped_msg",
  "data": {$data},
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $duration_ms
}
RESULT
}

is_integer() {
    case "$1" in
        ''|*[!0-9]*) return 1 ;;
        *) return 0 ;;
    esac
}

START_TIME=$(now_ms)

URL="${URL:-${CHECK_URL:-}}"
METHOD="${METHOD:-GET}"
EXPECTED_STATUS="${EXPECTED_STATUS:-200}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-30}"
FOGHORN_TIMEOUT_RAW="${FOGHORN_TIMEOUT:-}"
VERIFY_SSL_RAW="$(echo "${VERIFY_SSL:-true}" | tr '[:upper:]' '[:lower:]')"
FOLLOW_REDIRECTS_RAW="$(echo "${FOLLOW_REDIRECTS:-true}" | tr '[:upper:]' '[:lower:]')"
HEADERS="${HEADERS:-}"
REQUEST_BODY="${REQUEST_BODY:-}"
WARNING_THRESHOLD_MS="${WARNING_THRESHOLD_MS:-1000}"
CRITICAL_THRESHOLD_MS="${CRITICAL_THRESHOLD_MS:-5000}"
CONTENT_REGEX="${CONTENT_REGEX:-}"

if [ -z "$URL" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "URL is required (set URL or CHECK_URL env var)" "$DURATION" ""
    exit 1
fi

case "$URL" in
    http://*|https://*) ;;
    *)
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "fail" "invalid URL scheme: $URL (must be http:// or https://)" "$DURATION" ""
        exit 1
        ;;
esac

if ! is_integer "$TIMEOUT_SECONDS"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "invalid TIMEOUT_SECONDS: $TIMEOUT_SECONDS" "$DURATION" ""
    exit 1
fi

if [ -n "$FOGHORN_TIMEOUT_RAW" ]; then
    if is_integer "$FOGHORN_TIMEOUT_RAW"; then
        if [ "$FOGHORN_TIMEOUT_RAW" -lt "$TIMEOUT_SECONDS" ]; then
            TIMEOUT_SECONDS="$FOGHORN_TIMEOUT_RAW"
        fi
    fi
fi

SSL_FLAG=""
case "$VERIFY_SSL_RAW" in
    true|1|yes) ;;
    false|0|no) SSL_FLAG="-k" ;;
    *)
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "fail" "invalid VERIFY_SSL: must be true or false" "$DURATION" ""
        exit 1
        ;;
esac

REDIRECT_FLAG="-L"
case "$FOLLOW_REDIRECTS_RAW" in
    true|1|yes) ;;
    false|0|no) REDIRECT_FLAG="" ;;
    *)
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "fail" "invalid FOLLOW_REDIRECTS: must be true or false" "$DURATION" ""
        exit 1
        ;;
esac

if ! is_integer "$WARNING_THRESHOLD_MS"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "invalid WARNING_THRESHOLD_MS: $WARNING_THRESHOLD_MS" "$DURATION" ""
    exit 1
fi

if ! is_integer "$CRITICAL_THRESHOLD_MS"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "invalid CRITICAL_THRESHOLD_MS: $CRITICAL_THRESHOLD_MS" "$DURATION" ""
    exit 1
fi

HEADER_FLAGS=""
if [ -n "$HEADERS" ]; then
    HEADER_KEYS=$(printf '%s' "$HEADERS" | sed -n 's/.*"//;s/".*//;s/{//;s/}//;p' 2>/dev/null || true)
    if command -v jq >/dev/null 2>&1; then
        for key in $(printf '%s' "$HEADERS" | jq -r 'keys[]' 2>/dev/null); do
            value=$(printf '%s' "$HEADERS" | jq -r --arg k "$key" '.[$k]' 2>/dev/null)
            HEADER_FLAGS="$HEADER_FLAGS -H \"$key: $value\""
        done
    else
        PAIRS=$(printf '%s' "$HEADERS" | sed -e 's/[{}"]//g' -e 's/,/\n/g' 2>/dev/null || true)
        IFS_OLD="$IFS"
        IFS='
'
        for pair in $PAIRS; do
            key=$(printf '%s' "$pair" | sed 's/^[[:space:]]*//' | cut -d: -f1 | sed 's/[[:space:]]*$//')
            value=$(printf '%s' "$pair" | sed 's/^[^:]*:[[:space:]]*//')
            if [ -n "$key" ]; then
                HEADER_FLAGS="$HEADER_FLAGS -H \"$key: $value\""
            fi
        done
        IFS="$IFS_OLD"
    fi
fi

BODY_FLAG=""
if [ -n "$REQUEST_BODY" ]; then
    BODY_FLAG="-d \"$REQUEST_BODY\""
fi

BODY_FILE=$(mktemp)
HEADER_FILE=$(mktemp)
WRITE_OUT_FILE=$(mktemp)
trap 'rm -f "$BODY_FILE" "$HEADER_FILE" "$WRITE_OUT_FILE"' EXIT

CURL_ARGS="-sS -X $METHOD -o $BODY_FILE -D $HEADER_FILE --max-time $TIMEOUT_SECONDS -w '%{http_code}\\n%{time_total}\\n%{size_download}'"
if [ -n "$SSL_FLAG" ]; then
    CURL_ARGS="$CURL_ARGS $SSL_FLAG"
fi
if [ -n "$REDIRECT_FLAG" ]; then
    CURL_ARGS="$CURL_ARGS $REDIRECT_FLAG"
fi
CURL_ARGS="$CURL_ARGS $HEADER_FLAGS $BODY_FLAG \"$URL\""

eval curl $CURL_ARGS > "$WRITE_OUT_FILE" 2>/dev/null || true

HTTP_CODE=$(sed -n '1p' "$WRITE_OUT_FILE" 2>/dev/null || true)
TIME_TOTAL=$(sed -n '2p' "$WRITE_OUT_FILE" 2>/dev/null || true)
SIZE_DOWNLOAD=$(sed -n '3p' "$WRITE_OUT_FILE" 2>/dev/null || true)

case "$HTTP_CODE" in
    ''|*[!0-9]*)
        HTTP_CODE="000"
        ;;
esac

case "$TIME_TOTAL" in
    ''|*[!0-9.]*)
        TIME_TOTAL="0"
        ;;
esac

RESPONSE_TIME_MS=$(printf '%s' "$TIME_TOTAL" | awk '{printf "%.0f", $1 * 1000}')

case "$SIZE_DOWNLOAD" in
    ''|*[!0-9.]*)
        SIZE_DOWNLOAD="0"
        ;;
esac
RESPONSE_SIZE_BYTES=$(printf '%s' "$SIZE_DOWNLOAD" | awk '{printf "%.0f", $1}')

STATUS_TEXT=""
if [ -s "$HEADER_FILE" ]; then
    STATUS_LINE=$(head -1 "$HEADER_FILE" 2>/dev/null || true)
    STATUS_TEXT=$(printf '%s' "$STATUS_LINE" | sed 's/^[^ ]* [0-9]*[[:space:]]*//' | tr -d '\r\n' | sed 's/[[:space:]]*$//')
    if [ -z "$STATUS_TEXT" ]; then
        STATUS_TEXT="Unknown"
    fi
else
    STATUS_TEXT="Unknown"
fi

escaped_url="$(json_escape "$URL")"
escaped_method="$(json_escape "$METHOD")"
escaped_status_text="$(json_escape "$STATUS_TEXT")"

DATA_FIELDS="\"url\": \"$escaped_url\", \"method\": \"$escaped_method\", \"status_code\": $HTTP_CODE, \"status_text\": \"$escaped_status_text\", \"response_time_ms\": $RESPONSE_TIME_MS, \"response_size_bytes\": $RESPONSE_SIZE_BYTES"

STATUS_MATCH="false"

case "$EXPECTED_STATUS" in
    *-*)
        RANGE_START=$(printf '%s' "$EXPECTED_STATUS" | cut -d- -f1)
        RANGE_END=$(printf '%s' "$EXPECTED_STATUS" | cut -d- -f2)
        if is_integer "$RANGE_START" && is_integer "$RANGE_END" && \
           [ "$HTTP_CODE" -ge "$RANGE_START" ] && [ "$HTTP_CODE" -le "$RANGE_END" ]; then
            STATUS_MATCH="true"
        fi
        ;;
    *,*)
        IFS_OLD="$IFS"
        IFS=','
        for code in $EXPECTED_STATUS; do
            if [ "$code" = "$HTTP_CODE" ]; then
                STATUS_MATCH="true"
                break
            fi
        done
        IFS="$IFS_OLD"
        ;;
    *)
        if [ "$HTTP_CODE" = "$EXPECTED_STATUS" ]; then
            STATUS_MATCH="true"
        fi
        ;;
esac

if [ -n "$CONTENT_REGEX" ]; then
    if grep -qE "$CONTENT_REGEX" "$BODY_FILE" 2>/dev/null; then
        CONTENT_MATCH="true"
    else
        CONTENT_MATCH="false"
    fi
    DATA_FIELDS="$DATA_FIELDS, \"content_match\": $CONTENT_MATCH"
fi

DATA_FIELDS="$DATA_FIELDS"

END_TIME=$(now_ms)
DURATION=$((END_TIME - START_TIME))

if [ "$STATUS_MATCH" = "false" ]; then
    json_output "fail" "expected status $EXPECTED_STATUS, got $HTTP_CODE" "$DURATION" "$DATA_FIELDS"
    exit 1
fi

if [ "$HTTP_CODE" = "000" ]; then
    json_output "fail" "request failed or timed out" "$DURATION" "$DATA_FIELDS"
    exit 1
fi

if [ -n "$CONTENT_REGEX" ] && [ "$CONTENT_MATCH" = "false" ]; then
    json_output "fail" "content regex did not match" "$DURATION" "$DATA_FIELDS"
    exit 1
fi

STATUS="pass"
if [ "$RESPONSE_TIME_MS" -ge "$CRITICAL_THRESHOLD_MS" ]; then
    STATUS="fail"
elif [ "$RESPONSE_TIME_MS" -ge "$WARNING_THRESHOLD_MS" ]; then
    STATUS="warn"
fi

case "$STATUS" in
    pass) MESSAGE="HTTP check passed with status $HTTP_CODE in ${RESPONSE_TIME_MS}ms" ;;
    warn) MESSAGE="HTTP check warning: status $HTTP_CODE but response time ${RESPONSE_TIME_MS}ms >= ${WARNING_THRESHOLD_MS}ms" ;;
    fail) MESSAGE="HTTP check failed: response time ${RESPONSE_TIME_MS}ms >= ${CRITICAL_THRESHOLD_MS}ms critical threshold" ;;
esac

json_output "$STATUS" "$MESSAGE" "$DURATION" "$DATA_FIELDS"

if [ "$STATUS" = "pass" ]; then
    exit 0
fi

exit 1
