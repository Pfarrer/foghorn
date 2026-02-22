package secretstore

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseRef(t *testing.T) {
	key, ok := ParseRef("secret://smtp/password")
	if !ok {
		t.Fatalf("expected reference to parse")
	}
	if key != "smtp/password" {
		t.Fatalf("unexpected key: %s", key)
	}

	_, ok = ParseRef("SMTP_PASSWORD")
	if ok {
		t.Fatal("expected non-reference to fail parsing")
	}
}

func TestStoreCRUD(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.enc")

	store, err := New(path, []byte("test-master-key-thirty-two-bytes-long"))
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if err := store.Set("smtp/password", "abc123"); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if err := store.Set("imap/password", "xyz789"); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	keys, err := store.ListKeys()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	expectedKeys := []string{"imap/password", "smtp/password"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("unexpected keys: %#v", keys)
	}

	value, err := store.Resolve("secret://smtp/password")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if value != "abc123" {
		t.Fatalf("unexpected value: %s", value)
	}

	deleted, err := store.Delete("imap/password")
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if !deleted {
		t.Fatal("expected key to be deleted")
	}

	_, err = store.Resolve("secret://imap/password")
	if err == nil {
		t.Fatal("expected missing key resolve error")
	}
}

func TestStoreIsEncryptedAtRest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.enc")

	store, err := New(path, []byte("test-master-key-thirty-two-bytes-long"))
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	secretValue := "super-secret-value"
	if err := store.Set("smtp/password", secretValue); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read store: %v", err)
	}
	if string(payload) == "" {
		t.Fatal("expected non-empty store payload")
	}
	if contains(string(payload), secretValue) {
		t.Fatal("secret value should not appear in store payload")
	}
}

func TestMasterKeyFromEnv(t *testing.T) {
	old := os.Getenv("FOGHORN_SECRET_MASTER_KEY")
	defer os.Setenv("FOGHORN_SECRET_MASTER_KEY", old)

	os.Setenv("FOGHORN_SECRET_MASTER_KEY", "test-key-thirty-two-bytes-long-!")
	key, err := MasterKeyFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(key) != 32 {
		t.Fatalf("expected 32-byte key, got %d", len(key))
	}

	os.Unsetenv("FOGHORN_SECRET_MASTER_KEY")
	_, err = MasterKeyFromEnv()
	if err == nil {
		t.Fatal("expected missing env var error")
	}

	os.Setenv("FOGHORN_SECRET_MASTER_KEY", "short")
	_, err = MasterKeyFromEnv()
	if err == nil {
		t.Fatal("expected min length error")
	}
}

func TestSecretValueMaxSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.enc")

	store, err := New(path, []byte("test-master-key-thirty-two-bytes-long"))
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	largeValue := make([]byte, 64*1024+1)
	err = store.Set("test/large", string(largeValue))
	if err == nil {
		t.Fatal("expected max size error for large secret")
	}

	validValue := make([]byte, 64*1024)
	err = store.Set("test/valid", string(validValue))
	if err != nil {
		t.Fatalf("unexpected error for valid size: %v", err)
	}
}

func contains(s string, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
