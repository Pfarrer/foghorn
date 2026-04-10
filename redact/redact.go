package redact

import (
	"regexp"
	"slices"
	"strings"
)

var (
	authHeaderPattern  = regexp.MustCompile(`(?im)(authorization\s*[:=]\s*)([^\r\n]+)`)
	credentialPattern  = regexp.MustCompile(`(?im)("?(?:password|passwd|token|secret|api[_-]?key|authorization)"?\s*[:=]\s*)("[^"]*"|'[^']*'|[^\s,}]+)`)
	bearerTokenPattern = regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9\-._~+/]+=*`)
)

func Sanitize(input string, secrets []string) string {
	redacted := input

	uniqueSecrets := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		if !slices.Contains(uniqueSecrets, secret) {
			uniqueSecrets = append(uniqueSecrets, secret)
		}
	}
	for _, secret := range uniqueSecrets {
		redacted = strings.ReplaceAll(redacted, secret, "[REDACTED]")
	}

	redacted = authHeaderPattern.ReplaceAllString(redacted, "${1}[REDACTED]")
	redacted = credentialPattern.ReplaceAllString(redacted, "${1}[REDACTED]")
	redacted = bearerTokenPattern.ReplaceAllString(redacted, "Bearer [REDACTED]")

	return redacted
}
