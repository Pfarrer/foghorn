# Standard DNS Resolution Check Container

## Category
integration

## Description
Define a reusable Docker container that performs DNS resolution checks and validates DNS records

## Usage Steps
1. User pulls the standard DNS check image: `foghorn/dns-check:1.0.0`
2. Configure check with domain name and DNS server
3. Foghorn runs the container on schedule
4. Container performs DNS query and outputs results

## Implementation Notes
Container should accept these environment variables:
- `DOMAIN_NAME` (required): Domain name to resolve
- `DNS_SERVER`: DNS server to query (default: system default, usually 8.8.8.8)
- `RECORD_TYPE`: DNS record type - A, AAAA, MX, TXT, CNAME, NS, SOA (default: A)
- `EXPECTED_VALUE`: Expected record value for validation (optional)
- `WARNING_THRESHOLD_MS`: Query time threshold for "warn" status (default: 500ms)
- `CRITICAL_THRESHOLD_MS`: Query time threshold for "fail" status (default: 2000ms)

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "Resolved example.com to 93.184.216.34 (25ms)",
  "data": {
    "domain": "example.com",
    "dns_server": "8.8.8.8",
    "record_type": "A",
    "query_time_ms": 25,
    "answers": [
      {
        "type": "A",
        "value": "93.184.216.34",
        "ttl": 86400
      }
    ],
    "expected_value_match": true
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 32
}
```

Status determination:
- `pass`: Resolution successful, query time < WARNING_THRESHOLD_MS, record matches EXPECTED_VALUE if specified
- `warn`: Resolution successful, query time >= WARNING_THRESHOLD_MS but < CRITICAL_THRESHOLD_MS
- `fail`: Resolution failed (NXDOMAIN, SERVFAIL, timeout), query time >= CRITICAL_THRESHOLD_MS, or record doesn't match expected value
- `unknown`: Error executing query or invalid configuration

Multiple answers (e.g., for A records or MX records):
```json
"answers": [
  {"type": "A", "value": "93.184.216.34", "ttl": 86400},
  {"type": "A", "value": "93.184.216.35", "ttl": 86400}
]
```

For MX records, include priority:
```json
"answers": [
  {"type": "MX", "value": "mail.example.com", "priority": 10, "ttl": 3600}
]
```

Use Alpine Linux with `dig` or `nslookup`, or Go with `miekg/dns` library for better control

## Acceptance Criteria
- [ ] Container accepts DOMAIN_NAME and validates it's not empty
- [ ] Supports common DNS record types (A, AAAA, MX, TXT, CNAME, NS, SOA)
- [ ] Performs DNS query with specified server or system default
- [ ] Measures and reports query time
- [ ] Returns all answers in data object
- [ ] Optionally validates against expected value
- [ ] Handles NXDOMAIN and SERVFAIL responses appropriately
- [ ] Supports custom DNS server specification
- [ ] Returns appropriate status based on query success, thresholds, and expected value
- [ ] Example configurations for different record types
- [ ] Respects FOGHORN_TIMEOUT if provided

## Example Configuration
```yaml
name: "api-dns-check"
image: "foghorn/dns-check:1.0.0"
schedule:
  cron: "*/10 * * * *"
env:
  DOMAIN_NAME: "api.example.com"
  RECORD_TYPE: "A"
  DNS_SERVER: "8.8.8.8"
  WARNING_THRESHOLD_MS: "500"
  CRITICAL_THRESHOLD_MS: "2000"
timeout: "10s"
```

```yaml
name: "mx-record-check"
image: "foghorn/dns-check:1.0.0"
schedule:
  interval: "5m"
env:
  DOMAIN_NAME: "example.com"
  RECORD_TYPE: "MX"
  EXPECTED_VALUE: "mail.example.com"
timeout: "10s"
```

```yaml
name: "cdn-ipv6-check"
image: "foghorn/dns-check:1.0.0"
schedule:
  cron: "0 * * * *"
env:
  DOMAIN_NAME: "cdn.example.com"
  RECORD_TYPE: "AAAA"
  DNS_SERVER: "1.1.1.1"
timeout: "15s"
```

## Passes
false
