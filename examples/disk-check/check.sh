#!/bin/bash

set -e

START_TIME=$(date +%s%3N)

# Configuration
MOUNT_POINT="${MOUNT_POINT}"
WARNING_THRESHOLD_PERCENT="${WARNING_THRESHOLD_PERCENT:-80}"
CRITICAL_THRESHOLD_PERCENT="${CRITICAL_THRESHOLD_PERCENT:-90}"
WARNING_THRESHOLD_BYTES="${WARNING_THRESHOLD_BYTES:-}"
CRITICAL_THRESHOLD_BYTES="${CRITICAL_THRESHOLD_BYTES:-}"
INODE_WARNING_PERCENT="${INODE_WARNING_PERCENT:-85}"
INODE_CRITICAL_PERCENT="${INODE_CRITICAL_PERCENT:-95}"
CHECK_INODES="${CHECK_INODES:-true}"
FOGHORN_TIMEOUT="${FOGHORN_TIMEOUT:-10}"

# Validate mount point
if [ -z "$MOUNT_POINT" ]; then
    STATUS="fail"
    MESSAGE="MOUNT_POINT environment variable is required"
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

# Check if mount point exists
if [ ! -d "$MOUNT_POINT" ]; then
    END_TIME=$(date +%s%3N)
    DURATION=$((END_TIME - START_TIME))

    STATUS="unknown"
    MESSAGE="Mount point does not exist: $MOUNT_POINT"
    
    cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "mount_point": "$MOUNT_POINT"
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
EOF
    exit 1
fi

# Get disk space statistics using df
if ! DF_OUTPUT=$(df -B1 "$MOUNT_POINT" 2>&1); then
    END_TIME=$(date +%s%3N)
    DURATION=$((END_TIME - START_TIME))

    STATUS="unknown"
    MESSAGE="Failed to read filesystem information: $DF_OUTPUT"
    
    cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "mount_point": "$MOUNT_POINT"
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
EOF
    exit 1
fi

# Parse df output
# Skip header, get the line for our mount point
DISK_STATS=$(echo "$DF_OUTPUT" | awk "NR==2 {print \$2, \$3, \$4}")

TOTAL_BYTES=$(echo $DISK_STATS | awk '{print $1}')
USED_BYTES=$(echo $DISK_STATS | awk '{print $2}')
FREE_BYTES=$(echo $DISK_STATS | awk '{print $3}')

# Calculate usage percentage
if [ "$TOTAL_BYTES" -gt 0 ]; then
    USAGE_PERCENT=$((USED_BYTES * 100 / TOTAL_BYTES))
else
    USAGE_PERCENT=0
fi

# Get inode statistics
TOTAL_INODES=0
USED_INODES=0
FREE_INODES=0
INODE_USAGE_PERCENT=0

if [ "$CHECK_INODES" = "true" ]; then
    if ! INODE_OUTPUT=$(df -i "$MOUNT_POINT" 2>&1); then
        END_TIME=$(date +%s%3N)
        DURATION=$((END_TIME - START_TIME))

        STATUS="unknown"
        MESSAGE="Failed to read inode information: $INODE_OUTPUT"
        
        cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
  "data": {
    "mount_point": "$MOUNT_POINT",
    "total_bytes": $TOTAL_BYTES,
    "used_bytes": $USED_BYTES,
    "free_bytes": $FREE_BYTES,
    "usage_percent": $USAGE_PERCENT
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "duration_ms": $DURATION
}
EOF
        exit 1
    fi

    # Parse inode output - handle different df output formats
    # Try standard format: Filesystem    Inodes IUsed IFree IUse% Mounted on
    INODE_LINE=$(echo "$INODE_OUTPUT" | grep "$MOUNT_POINT" | head -1)
    if [ -n "$INODE_LINE" ]; then
        # Extract inode values from the line
        TOTAL_INODES=$(echo "$INODE_LINE" | awk '{print $2}')
        USED_INODES=$(echo "$INODE_LINE" | awk '{print $3}')
        FREE_INODES=$(echo "$INODE_LINE" | awk '{print $4}')

        # Convert to integers (remove commas if present)
        TOTAL_INODES=$(echo "$TOTAL_INODES" | tr -d ',')
        USED_INODES=$(echo "$USED_INODES" | tr -d ',')
        FREE_INODES=$(echo "$FREE_INODES" | tr -d ',')

        # Validate they are numbers
        if ! [[ "$TOTAL_INODES" =~ ^[0-9]+$ ]]; then TOTAL_INODES=0; fi
        if ! [[ "$USED_INODES" =~ ^[0-9]+$ ]]; then USED_INODES=0; fi
        if ! [[ "$FREE_INODES" =~ ^[0-9]+$ ]]; then FREE_INODES=0; fi

        if [ "$TOTAL_INODES" -gt 0 ]; then
            INODE_USAGE_PERCENT=$((USED_INODES * 100 / TOTAL_INODES))
        fi
    else
        # Inode information not available for this filesystem
        CHECK_INODES="false"
    fi
fi

END_TIME=$(date +%s%3N)
DURATION=$((END_TIME - START_TIME))

# Determine status
STATUS="pass"

# Critical thresholds
BYTES_CRITICAL=0
PERCENT_CRITICAL=0
INODE_CRITICAL=0

# Check byte-based critical threshold
if [ -n "$CRITICAL_THRESHOLD_BYTES" ]; then
    if [ "$USED_BYTES" -ge "$CRITICAL_THRESHOLD_BYTES" ]; then
        BYTES_CRITICAL=1
    fi
fi

# Check percent-based critical threshold
if [ "$USAGE_PERCENT" -ge "$CRITICAL_THRESHOLD_PERCENT" ]; then
    PERCENT_CRITICAL=1
fi

# Check inode critical threshold
if [ "$CHECK_INODES" = "true" ] && [ "$INODE_USAGE_PERCENT" -ge "$INODE_CRITICAL_PERCENT" ]; then
    INODE_CRITICAL=1
fi

# If byte thresholds are provided, use them exclusively
if [ -n "$CRITICAL_THRESHOLD_BYTES" ] || [ -n "$WARNING_THRESHOLD_BYTES" ]; then
    if [ "$BYTES_CRITICAL" -eq 1 ]; then
        STATUS="fail"
    elif [ -n "$WARNING_THRESHOLD_BYTES" ] && [ "$USED_BYTES" -ge "$WARNING_THRESHOLD_BYTES" ]; then
        STATUS="warn"
    fi
else
    # Use percentage thresholds
    if [ "$PERCENT_CRITICAL" -eq 1 ] || [ "$INODE_CRITICAL" -eq 1 ]; then
        STATUS="fail"
    else
        # Warning thresholds
        BYTES_WARNING=0
        PERCENT_WARNING=0
        INODE_WARNING=0

        if [ -n "$WARNING_THRESHOLD_BYTES" ] && [ "$USED_BYTES" -ge "$WARNING_THRESHOLD_BYTES" ]; then
            BYTES_WARNING=1
        fi

        if [ "$USAGE_PERCENT" -ge "$WARNING_THRESHOLD_PERCENT" ]; then
            PERCENT_WARNING=1
        fi

        if [ "$CHECK_INODES" = "true" ] && [ "$INODE_USAGE_PERCENT" -ge "$INODE_WARNING_PERCENT" ]; then
            INODE_WARNING=1
        fi

        if [ "$PERCENT_WARNING" -eq 1 ] || [ "$INODE_WARNING" -eq 1 ]; then
            STATUS="warn"
        fi
    fi
fi

# Build message
TOTAL_GB=$((TOTAL_BYTES / 1024 / 1024 / 1024))
FREE_GB=$((FREE_BYTES / 1024 / 1024 / 1024))
USED_GB=$((USED_BYTES / 1024 / 1024 / 1024))

if [ "$CHECK_INODES" = "true" ]; then
    MESSAGE="$MOUNT_POINT: ${USAGE_PERCENT}% used (${USED_GB}GB used, ${FREE_GB}GB free of ${TOTAL_GB}GB total), ${INODE_USAGE_PERCENT}% inodes used"
else
    MESSAGE="$MOUNT_POINT: ${USAGE_PERCENT}% used (${USED_GB}GB used, ${FREE_GB}GB free of ${TOTAL_GB}GB total)"
fi

# Output JSON
cat <<EOF
{
  "status": "$STATUS",
  "message": "$MESSAGE",
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