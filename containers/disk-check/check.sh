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

json_output() {
    status="$1"
    message="$2"
    duration_ms="$3"

    cat <<RESULT
{
  "status": "$status",
  "message": "$message",
  "data": {
    "mount_point": "$MOUNT_POINT",
    "total_bytes": $TOTAL_BYTES,
    "used_bytes": $USED_BYTES,
    "free_bytes": $FREE_BYTES,
    "usage_percent": $USAGE_PERCENT,
    "total_inodes": $TOTAL_INODES,
    "used_inodes": $USED_INODES,
    "free_inodes": $FREE_INODES,
    "inode_usage_percent": $INODE_USAGE_PERCENT
  },
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

MOUNT_POINT="${MOUNT_POINT:-}"
WARNING_THRESHOLD_PERCENT="${WARNING_THRESHOLD_PERCENT:-80}"
CRITICAL_THRESHOLD_PERCENT="${CRITICAL_THRESHOLD_PERCENT:-90}"
WARNING_THRESHOLD_BYTES="${WARNING_THRESHOLD_BYTES:-}"
CRITICAL_THRESHOLD_BYTES="${CRITICAL_THRESHOLD_BYTES:-}"
INODE_WARNING_PERCENT="${INODE_WARNING_PERCENT:-85}"
INODE_CRITICAL_PERCENT="${INODE_CRITICAL_PERCENT:-95}"
CHECK_INODES_RAW="$(echo "${CHECK_INODES:-true}" | tr '[:upper:]' '[:lower:]')"
TIMEOUT_SECONDS_RAW="${TIMEOUT_SECONDS:-10}"
FOGHORN_TIMEOUT_RAW="${FOGHORN_TIMEOUT:-}"

TOTAL_BYTES=0
USED_BYTES=0
FREE_BYTES=0
USAGE_PERCENT=0
TOTAL_INODES=0
USED_INODES=0
FREE_INODES=0
INODE_USAGE_PERCENT=0
INODE_NOTE="inode check disabled"

if [ -z "$MOUNT_POINT" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "MOUNT_POINT is required" "$DURATION"
    exit 1
fi

if [ ! -e "$MOUNT_POINT" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "mount point does not exist: $MOUNT_POINT" "$DURATION"
    exit 1
fi

TIMEOUT_SECONDS=$(parse_seconds "$TIMEOUT_SECONDS_RAW" || true)
if [ -z "$TIMEOUT_SECONDS" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "invalid TIMEOUT_SECONDS: $TIMEOUT_SECONDS_RAW" "$DURATION"
    exit 1
fi

if [ -n "$FOGHORN_TIMEOUT_RAW" ]; then
    FOGHORN_TIMEOUT_SECONDS=$(parse_seconds "$FOGHORN_TIMEOUT_RAW" || true)
    if [ -z "$FOGHORN_TIMEOUT_SECONDS" ]; then
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "unknown" "invalid FOGHORN_TIMEOUT: $FOGHORN_TIMEOUT_RAW" "$DURATION"
        exit 1
    fi

    if [ "$FOGHORN_TIMEOUT_SECONDS" -lt "$TIMEOUT_SECONDS" ]; then
        TIMEOUT_SECONDS="$FOGHORN_TIMEOUT_SECONDS"
    fi
fi

for value in \
    "$WARNING_THRESHOLD_PERCENT" \
    "$CRITICAL_THRESHOLD_PERCENT" \
    "$INODE_WARNING_PERCENT" \
    "$INODE_CRITICAL_PERCENT"
do
    if ! is_integer "$value"; then
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "unknown" "thresholds must be non-negative integers" "$DURATION"
        exit 1
    fi
done

if [ -n "$WARNING_THRESHOLD_BYTES" ] && ! is_integer "$WARNING_THRESHOLD_BYTES"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "WARNING_THRESHOLD_BYTES must be a non-negative integer" "$DURATION"
    exit 1
fi

if [ -n "$CRITICAL_THRESHOLD_BYTES" ] && ! is_integer "$CRITICAL_THRESHOLD_BYTES"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "CRITICAL_THRESHOLD_BYTES must be a non-negative integer" "$DURATION"
    exit 1
fi

if [ "$CRITICAL_THRESHOLD_PERCENT" -lt "$WARNING_THRESHOLD_PERCENT" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "CRITICAL_THRESHOLD_PERCENT must be >= WARNING_THRESHOLD_PERCENT" "$DURATION"
    exit 1
fi

if [ "$INODE_CRITICAL_PERCENT" -lt "$INODE_WARNING_PERCENT" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "INODE_CRITICAL_PERCENT must be >= INODE_WARNING_PERCENT" "$DURATION"
    exit 1
fi

if [ -n "$WARNING_THRESHOLD_BYTES" ] && [ -n "$CRITICAL_THRESHOLD_BYTES" ] && [ "$CRITICAL_THRESHOLD_BYTES" -lt "$WARNING_THRESHOLD_BYTES" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "CRITICAL_THRESHOLD_BYTES must be >= WARNING_THRESHOLD_BYTES" "$DURATION"
    exit 1
fi

DISK_LINE=$(run_with_timeout "$TIMEOUT_SECONDS" df -kP "$MOUNT_POINT" 2>/dev/null | awk 'NR==2' || true)
if [ -z "$DISK_LINE" ]; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "failed to read disk usage for $MOUNT_POINT" "$DURATION"
    exit 1
fi

TOTAL_KB=$(echo "$DISK_LINE" | awk '{print $2}')
USED_KB=$(echo "$DISK_LINE" | awk '{print $3}')
FREE_KB=$(echo "$DISK_LINE" | awk '{print $4}')
USAGE_PERCENT_RAW=$(echo "$DISK_LINE" | awk '{print $5}')
USAGE_PERCENT=${USAGE_PERCENT_RAW%%%}

if ! is_integer "$TOTAL_KB" || ! is_integer "$USED_KB" || ! is_integer "$FREE_KB" || ! is_integer "$USAGE_PERCENT"; then
    END_TIME=$(now_ms)
    DURATION=$((END_TIME - START_TIME))
    json_output "unknown" "invalid disk usage values returned by df" "$DURATION"
    exit 1
fi

TOTAL_BYTES=$((TOTAL_KB * 1024))
USED_BYTES=$((USED_KB * 1024))
FREE_BYTES=$((FREE_KB * 1024))

case "$CHECK_INODES_RAW" in
    true|1|yes)
        CHECK_INODES=true
        INODE_NOTE=""
        ;;
    false|0|no)
        CHECK_INODES=false
        ;;
    *)
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "unknown" "CHECK_INODES must be true or false" "$DURATION"
        exit 1
        ;;
esac

if [ "$CHECK_INODES" = "true" ]; then
    INODE_LINE=$(run_with_timeout "$TIMEOUT_SECONDS" df -Pi "$MOUNT_POINT" 2>/dev/null | awk 'NR==2' || true)
    if [ -z "$INODE_LINE" ]; then
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "unknown" "failed to read inode usage for $MOUNT_POINT" "$DURATION"
        exit 1
    fi

    TOTAL_INODES=$(echo "$INODE_LINE" | awk '{print $2}')
    USED_INODES=$(echo "$INODE_LINE" | awk '{print $3}')
    FREE_INODES=$(echo "$INODE_LINE" | awk '{print $4}')
    INODE_USAGE_PERCENT_RAW=$(echo "$INODE_LINE" | awk '{print $5}')
    INODE_USAGE_PERCENT=${INODE_USAGE_PERCENT_RAW%%%}

    if [ "$TOTAL_INODES" = "-" ] || [ "$USED_INODES" = "-" ] || [ "$FREE_INODES" = "-" ] || [ "$INODE_USAGE_PERCENT" = "-" ]; then
        CHECK_INODES=false
        TOTAL_INODES=0
        USED_INODES=0
        FREE_INODES=0
        INODE_USAGE_PERCENT=0
        INODE_NOTE="inode stats unavailable"
    elif ! is_integer "$TOTAL_INODES" || ! is_integer "$USED_INODES" || ! is_integer "$FREE_INODES" || ! is_integer "$INODE_USAGE_PERCENT"; then
        END_TIME=$(now_ms)
        DURATION=$((END_TIME - START_TIME))
        json_output "unknown" "invalid inode values returned by df" "$DURATION"
        exit 1
    fi
fi

STATUS="pass"

if [ "$USAGE_PERCENT" -ge "$CRITICAL_THRESHOLD_PERCENT" ]; then
    STATUS="fail"
elif [ "$USAGE_PERCENT" -ge "$WARNING_THRESHOLD_PERCENT" ]; then
    STATUS="warn"
fi

if [ -n "$CRITICAL_THRESHOLD_BYTES" ] && [ "$USED_BYTES" -ge "$CRITICAL_THRESHOLD_BYTES" ]; then
    STATUS="fail"
elif [ -n "$WARNING_THRESHOLD_BYTES" ] && [ "$USED_BYTES" -ge "$WARNING_THRESHOLD_BYTES" ] && [ "$STATUS" = "pass" ]; then
    STATUS="warn"
fi

if [ "$CHECK_INODES" = "true" ]; then
    if [ "$INODE_USAGE_PERCENT" -ge "$INODE_CRITICAL_PERCENT" ]; then
        STATUS="fail"
    elif [ "$INODE_USAGE_PERCENT" -ge "$INODE_WARNING_PERCENT" ] && [ "$STATUS" = "pass" ]; then
        STATUS="warn"
    fi
fi

if [ "$CHECK_INODES" = "true" ]; then
    MESSAGE="$MOUNT_POINT: ${USAGE_PERCENT}% used (${FREE_BYTES}B free of ${TOTAL_BYTES}B total), ${INODE_USAGE_PERCENT}% inodes used"
else
    MESSAGE="$MOUNT_POINT: ${USAGE_PERCENT}% used (${FREE_BYTES}B free of ${TOTAL_BYTES}B total), $INODE_NOTE"
fi

END_TIME=$(now_ms)
DURATION=$((END_TIME - START_TIME))
json_output "$STATUS" "$MESSAGE" "$DURATION"

if [ "$STATUS" = "pass" ]; then
    exit 0
fi

exit 1
