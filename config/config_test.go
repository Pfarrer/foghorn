package config

import (
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	data, err := os.ReadFile("../example.yaml")
	if err != nil {
		t.Fatalf("Failed to read example.yaml: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("example.yaml is empty")
	}

	content := string(data)
	expectedFields := []string{"name:", "image:", "schedule:", "evaluation:", "enabled:"}

	for _, field := range expectedFields {
		if !contains(content, field) {
			t.Errorf("example.yaml missing expected field: %s", field)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
