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

json_output() {
    status="$1"
    message="$2"
    duration_ms="$3"
    receive_time="$4"

    esc_message="$(json_escape "$message")"
    esc_mail_from="$(json_escape "$MAIL_FROM")"
    esc_mail_to="$(json_escape "$MAIL_TO")"
    esc_smtp_host="$(json_escape "$SMTP_HOST")"
    esc_receive_host="$(json_escape "$RECEIVE_HOST")"
    esc_receive_protocol="$(json_escape "$RECEIVE_PROTOCOL")"
    esc_message_id="$(json_escape "$MESSAGE_ID")"
    esc_correlation_id="$(json_escape "$CORRELATION_ID")"
    esc_send_time="$(json_escape "$SEND_TIME_ISO")"
    esc_receive_time="$(json_escape "$receive_time")"

    cat <<RESULT
{
  "status": "$status",
  "message": "$esc_message",
  "data": {
    "mail_from": "$esc_mail_from",
    "mail_to": "$esc_mail_to",
    "smtp_host": "$esc_smtp_host",
    "receive_host": "$esc_receive_host",
    "receive_protocol": "$esc_receive_protocol",
    "message_id": "$esc_message_id",
    "correlation_id": "$esc_correlation_id",
    "send_time": "$esc_send_time",
    "receive_time": "$esc_receive_time",
    "delivery_seconds": $DELIVERY_SECONDS,
    "deadline_seconds": $EFFECTIVE_DEADLINE_SECONDS
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $duration_ms
}
RESULT
}

emit_fail() {
    msg="$1"
    end_time="$(now_ms)"
    duration="$((end_time - START_TIME_MS))"
    json_output "fail" "$msg" "$duration" ""
    exit 1
}

emit_unknown() {
    msg="$1"
    end_time="$(now_ms)"
    duration="$((end_time - START_TIME_MS))"
    json_output "unknown" "$msg" "$duration" ""
    exit 1
}

START_TIME_MS="$(now_ms)"
START_TIME_EPOCH="$(date +%s)"

SMTP_HOST="${SMTP_HOST:-}"
SMTP_PORT_RAW="${SMTP_PORT:-587}"
SMTP_USERNAME="${SMTP_USERNAME:-}"
SMTP_PASSWORD="${SMTP_PASSWORD:-}"
SMTP_TLS_MODE="$(printf '%s' "${SMTP_TLS_MODE:-starttls}" | tr '[:upper:]' '[:lower:]')"
MAIL_FROM="${MAIL_FROM:-}"
MAIL_TO="${MAIL_TO:-}"
SUBJECT_PREFIX="${SUBJECT_PREFIX:-FOGHORN-CHECK}"
BODY_TEMPLATE="${BODY_TEMPLATE:-Foghorn mail check}"

RECEIVE_HOST="${RECEIVE_HOST:-}"
RECEIVE_PORT_RAW="${RECEIVE_PORT:-993}"
RECEIVE_USERNAME="${RECEIVE_USERNAME:-}"
RECEIVE_PASSWORD="${RECEIVE_PASSWORD:-}"
RECEIVE_TLS_RAW="$(printf '%s' "${RECEIVE_TLS:-true}" | tr '[:upper:]' '[:lower:]')"
RECEIVE_MAILBOX="${RECEIVE_MAILBOX:-INBOX}"

POLL_INTERVAL_SECONDS_RAW="${POLL_INTERVAL_SECONDS:-5}"
WARNING_THRESHOLD_SECONDS_RAW="${WARNING_THRESHOLD_SECONDS:-}"
DEADLINE_SECONDS_RAW="${DEADLINE_SECONDS:-}"
DELETE_AFTER_MATCH_RAW="$(printf '%s' "${DELETE_AFTER_MATCH:-false}" | tr '[:upper:]' '[:lower:]')"
FOGHORN_TIMEOUT_RAW="${FOGHORN_TIMEOUT:-}"

DELIVERY_SECONDS=0
EFFECTIVE_DEADLINE_SECONDS=0
MESSAGE_ID=""
CORRELATION_ID=""
SEND_TIME_ISO=""
RECEIVE_PROTOCOL="imap"

if [ -z "$SMTP_HOST" ] || [ -z "$SMTP_USERNAME" ] || [ -z "$SMTP_PASSWORD" ] || [ -z "$MAIL_FROM" ] || [ -z "$MAIL_TO" ]; then
    emit_unknown "SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, MAIL_FROM, and MAIL_TO are required"
fi

if [ -z "$RECEIVE_HOST" ] || [ -z "$RECEIVE_USERNAME" ] || [ -z "$RECEIVE_PASSWORD" ]; then
    emit_unknown "RECEIVE_HOST, RECEIVE_USERNAME, and RECEIVE_PASSWORD are required"
fi

if [ -z "$DEADLINE_SECONDS_RAW" ]; then
    emit_unknown "DEADLINE_SECONDS is required"
fi

if ! is_integer "$SMTP_PORT_RAW" || [ "$SMTP_PORT_RAW" -lt 1 ] || [ "$SMTP_PORT_RAW" -gt 65535 ]; then
    emit_unknown "SMTP_PORT must be an integer between 1 and 65535"
fi
SMTP_PORT="$SMTP_PORT_RAW"

if ! is_integer "$RECEIVE_PORT_RAW" || [ "$RECEIVE_PORT_RAW" -lt 1 ] || [ "$RECEIVE_PORT_RAW" -gt 65535 ]; then
    emit_unknown "RECEIVE_PORT must be an integer between 1 and 65535"
fi
RECEIVE_PORT="$RECEIVE_PORT_RAW"

POLL_INTERVAL_SECONDS="$(parse_seconds "$POLL_INTERVAL_SECONDS_RAW" || true)"
if [ -z "$POLL_INTERVAL_SECONDS" ]; then
    emit_unknown "invalid POLL_INTERVAL_SECONDS: $POLL_INTERVAL_SECONDS_RAW"
fi

DEADLINE_SECONDS="$(parse_seconds "$DEADLINE_SECONDS_RAW" || true)"
if [ -z "$DEADLINE_SECONDS" ]; then
    emit_unknown "invalid DEADLINE_SECONDS: $DEADLINE_SECONDS_RAW"
fi

if [ -n "$WARNING_THRESHOLD_SECONDS_RAW" ]; then
    WARNING_THRESHOLD_SECONDS="$(parse_seconds "$WARNING_THRESHOLD_SECONDS_RAW" || true)"
    if [ -z "$WARNING_THRESHOLD_SECONDS" ]; then
        emit_unknown "invalid WARNING_THRESHOLD_SECONDS: $WARNING_THRESHOLD_SECONDS_RAW"
    fi
else
    WARNING_THRESHOLD_SECONDS=""
fi

if [ -n "$FOGHORN_TIMEOUT_RAW" ]; then
    FOGHORN_TIMEOUT_SECONDS="$(parse_seconds "$FOGHORN_TIMEOUT_RAW" || true)"
    if [ -z "$FOGHORN_TIMEOUT_SECONDS" ]; then
        emit_unknown "invalid FOGHORN_TIMEOUT: $FOGHORN_TIMEOUT_RAW"
    fi
else
    FOGHORN_TIMEOUT_SECONDS=""
fi

case "$SMTP_TLS_MODE" in
    starttls|tls|none) ;;
    *)
        emit_unknown "SMTP_TLS_MODE must be one of: starttls, tls, none"
        ;;
esac

case "$RECEIVE_TLS_RAW" in
    true|1|yes)
        RECEIVE_TLS=true
        ;;
    false|0|no)
        RECEIVE_TLS=false
        ;;
    *)
        emit_unknown "RECEIVE_TLS must be true or false"
        ;;
esac

case "$DELETE_AFTER_MATCH_RAW" in
    true|1|yes)
        DELETE_AFTER_MATCH=true
        ;;
    false|0|no)
        DELETE_AFTER_MATCH=false
        ;;
    *)
        emit_unknown "DELETE_AFTER_MATCH must be true or false"
        ;;
esac

EFFECTIVE_DEADLINE_SECONDS="$DEADLINE_SECONDS"
if [ -n "$FOGHORN_TIMEOUT_SECONDS" ] && [ "$FOGHORN_TIMEOUT_SECONDS" -lt "$EFFECTIVE_DEADLINE_SECONDS" ]; then
    EFFECTIVE_DEADLINE_SECONDS="$FOGHORN_TIMEOUT_SECONDS"
fi

END_DEADLINE_EPOCH="$((START_TIME_EPOCH + EFFECTIVE_DEADLINE_SECONDS))"

RANDOM_HEX="$(od -An -N4 -tx1 /dev/urandom 2>/dev/null | tr -d ' \n' || true)"
if [ -z "$RANDOM_HEX" ]; then
    RANDOM_HEX="$(date +%s)"
fi
CORRELATION_ID="foghorn-$(date -u +%Y%m%dT%H%M%SZ)-$RANDOM_HEX"
MESSAGE_ID="<${CORRELATION_ID}@foghorn.local>"
SUBJECT="${SUBJECT_PREFIX} ${CORRELATION_ID}"
SEND_TIME_ISO="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
MAIL_DATE="$(LC_ALL=C date -u +"%a, %d %b %Y %H:%M:%S +0000")"

MESSAGE_FILE="$(mktemp)"
cleanup() {
    rm -f "$MESSAGE_FILE"
}
trap cleanup EXIT

cat >"$MESSAGE_FILE" <<MAIL
From: ${MAIL_FROM}
To: ${MAIL_TO}
Subject: ${SUBJECT}
Date: ${MAIL_DATE}
Message-ID: ${MESSAGE_ID}
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8

${BODY_TEMPLATE}

correlation_id=${CORRELATION_ID}
sent_at=${SEND_TIME_ISO}
MAIL

case "$SMTP_TLS_MODE" in
    tls)
        SMTP_URL="smtps://${SMTP_HOST}:${SMTP_PORT}"
        SMTP_SSL_FLAG=""
        ;;
    starttls)
        SMTP_URL="smtp://${SMTP_HOST}:${SMTP_PORT}"
        SMTP_SSL_FLAG="--ssl-reqd"
        ;;
    none)
        SMTP_URL="smtp://${SMTP_HOST}:${SMTP_PORT}"
        SMTP_SSL_FLAG=""
        ;;
esac

SECONDS_LEFT="$((END_DEADLINE_EPOCH - $(date +%s)))"
if [ "$SECONDS_LEFT" -le 0 ]; then
    emit_fail "deadline exceeded before send started"
fi

set +e
if [ -n "$SMTP_SSL_FLAG" ]; then
    curl --silent --show-error --max-time "$SECONDS_LEFT" --connect-timeout 10 \
        --url "$SMTP_URL" \
        --user "${SMTP_USERNAME}:${SMTP_PASSWORD}" \
        --mail-from "$MAIL_FROM" \
        --mail-rcpt "$MAIL_TO" \
        --upload-file "$MESSAGE_FILE" \
        $SMTP_SSL_FLAG >/dev/null 2>&1
else
    curl --silent --show-error --max-time "$SECONDS_LEFT" --connect-timeout 10 \
        --url "$SMTP_URL" \
        --user "${SMTP_USERNAME}:${SMTP_PASSWORD}" \
        --mail-from "$MAIL_FROM" \
        --mail-rcpt "$MAIL_TO" \
        --upload-file "$MESSAGE_FILE" >/dev/null 2>&1
fi
SMTP_EXIT=$?
set -e
if [ "$SMTP_EXIT" -ne 0 ]; then
    emit_fail "failed to send mail via SMTP"
fi

SEND_EPOCH="$(date +%s)"

if [ "$RECEIVE_TLS" = "true" ]; then
    IMAP_URL="imaps://${RECEIVE_HOST}:${RECEIVE_PORT}/${RECEIVE_MAILBOX}"
else
    IMAP_URL="imap://${RECEIVE_HOST}:${RECEIVE_PORT}/${RECEIVE_MAILBOX}"
fi

FOUND_UID=""
FOUND_RECEIVE_TIME=""

while :; do
    NOW_EPOCH="$(date +%s)"
    SECONDS_LEFT="$((END_DEADLINE_EPOCH - NOW_EPOCH))"
    if [ "$SECONDS_LEFT" -le 0 ]; then
        emit_fail "message did not arrive within ${EFFECTIVE_DEADLINE_SECONDS}s"
    fi

    SEARCH_TIMEOUT="$SECONDS_LEFT"
    if [ "$SEARCH_TIMEOUT" -gt 15 ]; then
        SEARCH_TIMEOUT=15
    fi

    set +e
    SEARCH_OUTPUT="$(curl --silent --show-error --max-time "$SEARCH_TIMEOUT" --connect-timeout 10 \
        --url "$IMAP_URL" \
        --user "${RECEIVE_USERNAME}:${RECEIVE_PASSWORD}" \
        -X "UID SEARCH HEADER SUBJECT \"${CORRELATION_ID}\"" 2>/dev/null)"
    SEARCH_EXIT=$?
    set -e
    if [ "$SEARCH_EXIT" -ne 0 ]; then
        emit_fail "failed to query IMAP mailbox"
    fi

    FOUND_UID="$(printf '%s\n' "$SEARCH_OUTPUT" \
        | sed -n 's/^\* SEARCH[[:space:]]*//p' \
        | tr ' ' '\n' \
        | sed '/^$/d' \
        | tail -n 1)"

    if [ -n "$FOUND_UID" ]; then
        FETCH_TIMEOUT="$SECONDS_LEFT"
        if [ "$FETCH_TIMEOUT" -gt 15 ]; then
            FETCH_TIMEOUT=15
        fi

        PARSED_RECEIVE_EPOCH=""
        set +e
        FETCH_OUTPUT="$(curl --silent --show-error --max-time "$FETCH_TIMEOUT" --connect-timeout 10 \
            --url "$IMAP_URL" \
            --user "${RECEIVE_USERNAME}:${RECEIVE_PASSWORD}" \
            -X "UID FETCH ${FOUND_UID} (INTERNALDATE)" 2>/dev/null)"
        FETCH_EXIT=$?
        set -e
        if [ "$FETCH_EXIT" -eq 0 ]; then
            INTERNAL_DATE="$(printf '%s\n' "$FETCH_OUTPUT" | sed -n 's/.*INTERNALDATE \"\([^\"]*\)\".*/\1/p' | head -n 1)"
            if [ -n "$INTERNAL_DATE" ]; then
                PARSED_RECEIVE_EPOCH="$(date -u -d "$INTERNAL_DATE" +%s 2>/dev/null || true)"
            fi
        fi

        if [ -n "$PARSED_RECEIVE_EPOCH" ] && [ "$PARSED_RECEIVE_EPOCH" -lt "$START_TIME_EPOCH" ]; then
            FOUND_UID=""
        else
            if [ -z "$PARSED_RECEIVE_EPOCH" ]; then
                PARSED_RECEIVE_EPOCH="$(date +%s)"
            fi
            DELIVERY_SECONDS="$((PARSED_RECEIVE_EPOCH - SEND_EPOCH))"
            if [ "$DELIVERY_SECONDS" -lt 0 ]; then
                DELIVERY_SECONDS=0
            fi
            FOUND_RECEIVE_TIME="$(date -u -d "@${PARSED_RECEIVE_EPOCH}" +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u +"%Y-%m-%dT%H:%M:%SZ")"

            if [ "$DELETE_AFTER_MATCH" = "true" ]; then
                set +e
                curl --silent --show-error --max-time 10 --connect-timeout 5 \
                    --url "$IMAP_URL" \
                    --user "${RECEIVE_USERNAME}:${RECEIVE_PASSWORD}" \
                    -X "UID STORE ${FOUND_UID} +FLAGS.SILENT (\\Deleted)" >/dev/null 2>&1
                curl --silent --show-error --max-time 10 --connect-timeout 5 \
                    --url "$IMAP_URL" \
                    --user "${RECEIVE_USERNAME}:${RECEIVE_PASSWORD}" \
                    -X "EXPUNGE" >/dev/null 2>&1
                set -e
            fi

            break
        fi
    fi

    NOW_EPOCH="$(date +%s)"
    SECONDS_LEFT="$((END_DEADLINE_EPOCH - NOW_EPOCH))"
    if [ "$SECONDS_LEFT" -le 0 ]; then
        emit_fail "message did not arrive within ${EFFECTIVE_DEADLINE_SECONDS}s"
    fi

    SLEEP_FOR="$POLL_INTERVAL_SECONDS"
    if [ "$SLEEP_FOR" -gt "$SECONDS_LEFT" ]; then
        SLEEP_FOR="$SECONDS_LEFT"
    fi
    sleep "$SLEEP_FOR"
done

END_TIME_MS="$(now_ms)"
DURATION_MS="$((END_TIME_MS - START_TIME_MS))"

if [ -n "$WARNING_THRESHOLD_SECONDS" ] && [ "$DELIVERY_SECONDS" -gt "$WARNING_THRESHOLD_SECONDS" ]; then
    json_output "warn" "mail delivered in ${DELIVERY_SECONDS}s (warning threshold ${WARNING_THRESHOLD_SECONDS}s)" "$DURATION_MS" "$FOUND_RECEIVE_TIME"
    exit 0
fi

json_output "pass" "mail delivered in ${DELIVERY_SECONDS}s" "$DURATION_MS" "$FOUND_RECEIVE_TIME"
exit 0
