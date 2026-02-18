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

is_integer() {
    case "$1" in
        ''|*[!0-9]*) return 1 ;;
        *) return 0 ;;
    esac
}

parse_seconds() {
    value="$1"
    case "$value" in
        ''|*[!0-9smh]*) return 1 ;;
        *s)
            seconds="${value%s}"
            ;;
        *m)
            minutes="${value%m}"
            is_integer "$minutes" || return 1
            seconds=$((minutes * 60))
            ;;
        *h)
            hours="${value%h}"
            is_integer "$hours" || return 1
            seconds=$((hours * 3600))
            ;;
        *)
            seconds="$value"
            ;;
    esac

    is_integer "$seconds" || return 1
    [ "$seconds" -gt 0 ] || return 1

    echo "$seconds"
}

json_escape() {
    printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

json_bool() {
    case "$1" in
        true) echo "true" ;;
        false) echo "false" ;;
        *) echo "false" ;;
    esac
}

json_output() {
    status="$1"
    message="$2"
    duration_ms="$3"

    esc_message="$(json_escape "$message")"
    esc_host="$(json_escape "$HOST")"
    esc_sni="$(json_escape "$SNI")"
    esc_tls_version="$(json_escape "$TLS_VERSION")"
    esc_cipher="$(json_escape "$CIPHER")"
    esc_subject="$(json_escape "$SUBJECT")"
    esc_issuer="$(json_escape "$ISSUER")"
    esc_not_before="$(json_escape "$NOT_BEFORE")"
    esc_not_after="$(json_escape "$NOT_AFTER")"

    cat <<RESULT
{
  "status": "$status",
  "message": "$esc_message",
  "data": {
    "host": "$esc_host",
    "port": $PORT,
    "sni": "$esc_sni",
    "tls_version": "$esc_tls_version",
    "cipher": "$esc_cipher",
    "subject": "$esc_subject",
    "issuer": "$esc_issuer",
    "not_before": "$esc_not_before",
    "not_after": "$esc_not_after",
    "days_remaining": $DAYS_REMAINING,
    "trusted": $(json_bool "$TRUSTED"),
    "hostname_match": $(json_bool "$HOSTNAME_MATCH")
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $duration_ms
}
RESULT
}

run_with_timeout() {
    seconds="$1"
    shift

    if command -v timeout >/dev/null 2>&1; then
        timeout "${seconds}s" "$@"
        return $?
    fi

    "$@"
}

START_TIME=$(now_ms)

HOST="${HOST:-}"
PORT="${PORT:-}"
SNI="${SNI:-$HOST}"
MIN_TLS_VERSION_RAW="${MIN_TLS_VERSION:-1.2}"
CA_BUNDLE_PATH="${CA_BUNDLE_PATH:-}"
VERIFY_HOSTNAME_RAW="$(echo "${VERIFY_HOSTNAME:-true}" | tr '[:upper:]' '[:lower:]')"
WARNING_DAYS_RAW="${WARNING_DAYS:-30}"
TIMEOUT_SECONDS_RAW="${TIMEOUT_SECONDS:-10}"
FOGHORN_TIMEOUT_RAW="${FOGHORN_TIMEOUT:-}"

TLS_VERSION=""
CIPHER=""
SUBJECT=""
ISSUER=""
NOT_BEFORE=""
NOT_AFTER=""
DAYS_REMAINING=0
TRUSTED=false
HOSTNAME_MATCH=false

if [ -z "$HOST" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "HOST is required" "$DURATION"
    exit 1
fi

if [ -z "$PORT" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "PORT is required" "$DURATION"
    exit 1
fi

if ! is_integer "$PORT" || [ "$PORT" -lt 1 ] || [ "$PORT" -gt 65535 ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "PORT must be an integer between 1 and 65535" "$DURATION"
    exit 1
fi

if [ -n "$CA_BUNDLE_PATH" ] && [ ! -f "$CA_BUNDLE_PATH" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "CA_BUNDLE_PATH not found: $CA_BUNDLE_PATH" "$DURATION"
    exit 1
fi

WARNING_DAYS="$WARNING_DAYS_RAW"
if ! is_integer "$WARNING_DAYS"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "WARNING_DAYS must be a non-negative integer" "$DURATION"
    exit 1
fi

TIMEOUT_SECONDS=$(parse_seconds "$TIMEOUT_SECONDS_RAW" || true)
if [ -z "$TIMEOUT_SECONDS" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "invalid TIMEOUT_SECONDS: $TIMEOUT_SECONDS_RAW" "$DURATION"
    exit 1
fi

if [ -n "$FOGHORN_TIMEOUT_RAW" ]; then
    FOGHORN_TIMEOUT_SECONDS=$(parse_seconds "$FOGHORN_TIMEOUT_RAW" || true)
    if [ -z "$FOGHORN_TIMEOUT_SECONDS" ]; then
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "fail" "invalid FOGHORN_TIMEOUT: $FOGHORN_TIMEOUT_RAW" "$DURATION"
        exit 1
    fi

    if [ "$FOGHORN_TIMEOUT_SECONDS" -lt "$TIMEOUT_SECONDS" ]; then
        TIMEOUT_SECONDS="$FOGHORN_TIMEOUT_SECONDS"
    fi
fi

case "$VERIFY_HOSTNAME_RAW" in
    true|1|yes)
        VERIFY_HOSTNAME=true
        ;;
    false|0|no)
        VERIFY_HOSTNAME=false
        ;;
    *)
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "fail" "VERIFY_HOSTNAME must be true or false" "$DURATION"
        exit 1
        ;;
esac

case "$MIN_TLS_VERSION_RAW" in
    1.0|TLSv1|TLS1.0)
        MIN_PROTOCOL="TLSv1"
        ;;
    1.1|TLSv1.1|TLS1.1)
        MIN_PROTOCOL="TLSv1.1"
        ;;
    1.2|TLSv1.2|TLS1.2)
        MIN_PROTOCOL="TLSv1.2"
        ;;
    1.3|TLSv1.3|TLS1.3)
        MIN_PROTOCOL="TLSv1.3"
        ;;
    *)
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "fail" "MIN_TLS_VERSION must be one of 1.0, 1.1, 1.2, or 1.3" "$DURATION"
        exit 1
        ;;
esac

TLS_OUTPUT_FILE="$(mktemp)"
CERT_FILE="$(mktemp)"
cleanup() {
    rm -f "$TLS_OUTPUT_FILE" "$CERT_FILE"
}
trap cleanup EXIT

set -- \
    openssl s_client \
    -connect "$HOST:$PORT" \
    -servername "$SNI" \
    -showcerts \
    -verify_return_error \
    -min_protocol "$MIN_PROTOCOL"

if [ -n "$CA_BUNDLE_PATH" ]; then
    set -- "$@" -CAfile "$CA_BUNDLE_PATH"
fi

if [ "$VERIFY_HOSTNAME" = "true" ]; then
    set -- "$@" -verify_hostname "$HOST"
fi

set +e
run_with_timeout "$TIMEOUT_SECONDS" "$@" </dev/null >"$TLS_OUTPUT_FILE" 2>&1
SCLIENT_EXIT=$?
set -e

if [ "$SCLIENT_EXIT" -eq 124 ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "TLS check timed out after ${TIMEOUT_SECONDS}s for ${HOST}:${PORT}" "$DURATION"
    exit 1
fi

if [ "$SCLIENT_EXIT" -ne 0 ]; then
    ERR_LINE="$(tail -n 1 "$TLS_OUTPUT_FILE" 2>/dev/null || true)"
    if [ -z "$ERR_LINE" ]; then
        ERR_LINE="TLS handshake failed"
    fi
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "fail" "$ERR_LINE" "$DURATION"
    exit 1
fi

awk '
BEGIN {in_cert=0; printed=0}
/-----BEGIN CERTIFICATE-----/ {if (printed==0) {in_cert=1}}
{if (in_cert==1) print}
/-----END CERTIFICATE-----/ {if (in_cert==1) {printed=1; exit}}
' "$TLS_OUTPUT_FILE" >"$CERT_FILE"

if [ ! -s "$CERT_FILE" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "unable to parse peer certificate from openssl output" "$DURATION"
    exit 1
fi

TRUSTED=false
if grep -q "Verify return code: 0 (ok)" "$TLS_OUTPUT_FILE"; then
    TRUSTED=true
fi

TLS_VERSION=$(awk -F':' '/Protocol[[:space:]]*:/ {gsub(/^[[:space:]]+/, "", $2); print $2; exit}' "$TLS_OUTPUT_FILE")
if [ -z "$TLS_VERSION" ]; then
    TLS_VERSION=$(awk -F',' '/New, TLSv/ {gsub(/^[[:space:]]+/, "", $2); gsub(/[[:space:]]+$/, "", $2); print $2; exit}' "$TLS_OUTPUT_FILE")
fi

CIPHER=$(awk -F':' '/Cipher[[:space:]]*:/ {gsub(/^[[:space:]]+/, "", $2); print $2; exit}' "$TLS_OUTPUT_FILE")
if [ -z "$CIPHER" ]; then
    CIPHER=$(awk -F'Cipher is ' '/Cipher is / {print $2; exit}' "$TLS_OUTPUT_FILE")
fi

SUBJECT_RAW=$(openssl x509 -in "$CERT_FILE" -noout -subject 2>/dev/null || true)
ISSUER_RAW=$(openssl x509 -in "$CERT_FILE" -noout -issuer 2>/dev/null || true)
STARTDATE_RAW=$(openssl x509 -in "$CERT_FILE" -noout -startdate 2>/dev/null || true)
ENDDATE_RAW=$(openssl x509 -in "$CERT_FILE" -noout -enddate 2>/dev/null || true)

SUBJECT="${SUBJECT_RAW#subject=}"
ISSUER="${ISSUER_RAW#issuer=}"
STARTDATE_VALUE="${STARTDATE_RAW#notBefore=}"
ENDDATE_VALUE="${ENDDATE_RAW#notAfter=}"

if [ -z "$STARTDATE_VALUE" ] || [ -z "$ENDDATE_VALUE" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "unable to parse certificate validity dates" "$DURATION"
    exit 1
fi

NOT_BEFORE="$(date -u -d "$STARTDATE_VALUE" +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || true)"
NOT_AFTER="$(date -u -d "$ENDDATE_VALUE" +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || true)"

if [ -z "$NOT_BEFORE" ] || [ -z "$NOT_AFTER" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "unable to normalize certificate validity timestamps" "$DURATION"
    exit 1
fi

NOT_AFTER_EPOCH=$(date -u -d "$ENDDATE_VALUE" +%s 2>/dev/null || true)
NOW_EPOCH=$(date -u +%s)

if [ -z "$NOT_AFTER_EPOCH" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "unable to calculate certificate expiry" "$DURATION"
    exit 1
fi

if [ "$NOT_AFTER_EPOCH" -le "$NOW_EPOCH" ]; then
    DAYS_REMAINING=0
else
    DAYS_REMAINING=$(( (NOT_AFTER_EPOCH - NOW_EPOCH) / 86400 ))
fi

if openssl x509 -in "$CERT_FILE" -noout -checkhost "$HOST" >/dev/null 2>&1; then
    HOSTNAME_MATCH=true
else
    HOSTNAME_MATCH=false
fi

STATUS="pass"
MESSAGE="TLS ok for ${HOST}:${PORT} (expires in ${DAYS_REMAINING}d)"

if [ "$TRUSTED" != "true" ]; then
    STATUS="fail"
    MESSAGE="TLS certificate is not trusted for ${HOST}:${PORT}"
fi

if [ "$VERIFY_HOSTNAME" = "true" ] && [ "$HOSTNAME_MATCH" != "true" ]; then
    STATUS="fail"
    MESSAGE="TLS certificate hostname mismatch for ${HOST}:${PORT}"
fi

if [ "$NOT_AFTER_EPOCH" -le "$NOW_EPOCH" ]; then
    STATUS="fail"
    MESSAGE="TLS certificate expired for ${HOST}:${PORT}"
fi

if [ "$STATUS" = "pass" ] && [ "$DAYS_REMAINING" -le "$WARNING_DAYS" ]; then
    STATUS="warn"
    MESSAGE="TLS ok for ${HOST}:${PORT} (expires in ${DAYS_REMAINING}d)"
fi

END_TIME=$(now_ms)
DURATION=$((END_TIME - START_TIME))
json_output "$STATUS" "$MESSAGE" "$DURATION"

if [ "$STATUS" = "pass" ] || [ "$STATUS" = "warn" ]; then
    exit 0
fi

exit 1
