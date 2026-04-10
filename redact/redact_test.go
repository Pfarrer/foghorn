package redact

import (
	"strings"
	"testing"
)

func TestSanitizeSecretValues(t *testing.T) {
	input := "The password is super-secret-password and it's important"
	result := Sanitize(input, []string{"super-secret-password"})
	if strings.Contains(result, "super-secret-password") {
		t.Fatalf("secret value should be redacted, got: %s", result)
	}
	if !strings.Contains(result, "[REDACTED]") {
		t.Fatalf("expected [REDACTED] marker, got: %s", result)
	}
}

func TestSanitizeAuthHeader(t *testing.T) {
	input := "Authorization: Bearer token123"
	result := Sanitize(input, nil)
	if strings.Contains(result, "token123") {
		t.Fatalf("bearer token should be redacted, got: %s", result)
	}
	if !strings.Contains(result, "[REDACTED]") {
		t.Fatalf("expected [REDACTED] marker, got: %s", result)
	}
}

func TestSanitizePasswordPattern(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"password= value", "password=secret123"},
		{"token= value", "token=abc"},
		{`"password": "value"`, `"password": "secret123"`},
		{`"token": "value"`, `"token": "abc"`},
		{"secret= value", "secret=mysecret"},
		{"api_key= value", "api_key=key123"},
		{"api-key= value", "api-key=key123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sanitize(tt.input, nil)
			if strings.Contains(result, "secret123") || strings.Contains(result, "abc") || strings.Contains(result, "mysecret") || strings.Contains(result, "key123") {
				t.Fatalf("credential value should be redacted, got: %s", result)
			}
			if !strings.Contains(result, "[REDACTED]") {
				t.Fatalf("expected [REDACTED] marker, got: %s", result)
			}
		})
	}
}

func TestSanitizeMultipleSecrets(t *testing.T) {
	input := "user=admin pass=hunter2 token=abc123"
	result := Sanitize(input, []string{"hunter2", "abc123"})
	if strings.Contains(result, "hunter2") {
		t.Fatalf("first secret should be redacted, got: %s", result)
	}
	if strings.Contains(result, "abc123") {
		t.Fatalf("second secret should be redacted, got: %s", result)
	}
}

func TestSanitizeEmptySecrets(t *testing.T) {
	input := "no secrets here"
	result := Sanitize(input, nil)
	if result != input {
		t.Fatalf("expected unchanged input, got: %s", result)
	}
}

func TestSanitizeEmptyInput(t *testing.T) {
	result := Sanitize("", []string{"secret"})
	if result != "" {
		t.Fatalf("expected empty string, got: %s", result)
	}
}

func TestSanitizeDuplicateSecrets(t *testing.T) {
	input := "myval occurs: myval myval"
	result := Sanitize(input, []string{"myval", "myval"})
	count := strings.Count(result, "[REDACTED]")
	if count != 3 {
		t.Fatalf("expected 3 redaction markers, got %d: %s", count, result)
	}
}

func TestSanitizeEmptySecretValue(t *testing.T) {
	input := "some text"
	result := Sanitize(input, []string{"", ""})
	if result != input {
		t.Fatalf("empty secrets should not change input, got: %s", result)
	}
}
