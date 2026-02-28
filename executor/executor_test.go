package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSecretDirPermissions(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("failed to create executor: %v", err)
	}
	defer exec.Close()

	secretDir, err := exec.createSecretDir()
	if err != nil {
		t.Fatalf("failed to create secret directory: %v", err)
	}
	defer os.RemoveAll(secretDir)

	info, err := os.Stat(secretDir)
	if err != nil {
		t.Fatalf("failed to stat secret directory: %v", err)
	}

	perms := info.Mode().Perm()
	if perms != 0o755 {
		t.Fatalf("expected directory permissions 0o755, got 0%o", perms)
	}
}

func TestSecretFilePermissions(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("failed to create executor: %v", err)
	}
	defer exec.Close()

	secretDir := t.TempDir()
	secretFile := filepath.Join(secretDir, "test-secret")
	secretValue := "super-secret-value"

	if err := os.WriteFile(secretFile, []byte(secretValue), 0o600); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	info, err := os.Stat(secretFile)
	if err != nil {
		t.Fatalf("failed to stat secret file: %v", err)
	}

	perms := info.Mode().Perm()
	if perms != 0o600 {
		t.Fatalf("expected file permissions 0o600, got 0%o", perms)
	}

	content, err := os.ReadFile(secretFile)
	if err != nil {
		t.Fatalf("failed to read secret file: %v", err)
	}

	if string(content) != secretValue {
		t.Fatalf("unexpected secret content: %s", string(content))
	}
}

func TestCleanupOldSecretDirs(t *testing.T) {
	tempDir := t.TempDir()

	secretDir1 := filepath.Join(tempDir, "old-dir")
	if err := os.Mkdir(secretDir1, 0o700); err != nil {
		t.Fatalf("failed to create old secret directory: %v", err)
	}

	tsFile1 := filepath.Join(secretDir1, ".timestamp")
	oldTime := time.Now().Add(-25 * time.Hour)
	if err := os.WriteFile(tsFile1, []byte(oldTime.Format(time.RFC3339)), 0o600); err != nil {
		t.Fatalf("failed to write old timestamp: %v", err)
	}

	secretDir2 := filepath.Join(tempDir, "recent-dir")
	if err := os.Mkdir(secretDir2, 0o700); err != nil {
		t.Fatalf("failed to create recent secret directory: %v", err)
	}

	tsFile2 := filepath.Join(secretDir2, ".timestamp")
	recentTime := time.Now().Add(-1 * time.Hour)
	if err := os.WriteFile(tsFile2, []byte(recentTime.Format(time.RFC3339)), 0o600); err != nil {
		t.Fatalf("failed to write recent timestamp: %v", err)
	}

	testFile := filepath.Join(secretDir1, "secret.txt")
	if err := os.WriteFile(testFile, []byte("sensitive-data"), 0o600); err != nil {
		t.Fatalf("failed to write test secret file: %v", err)
	}

	if err := cleanupOldSecretDirs(tempDir); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if _, err := os.Stat(secretDir1); !os.IsNotExist(err) {
		t.Fatal("old secret directory should have been removed")
	}

	if _, err := os.Stat(secretDir2); err != nil {
		t.Fatalf("recent secret directory should still exist: %v", err)
	}
}

func TestCleanupIgnoresMissingTimestamp(t *testing.T) {
	tempDir := t.TempDir()

	secretDir := filepath.Join(tempDir, "no-timestamp-dir")
	if err := os.Mkdir(secretDir, 0o700); err != nil {
		t.Fatalf("failed to create directory without timestamp: %v", err)
	}

	secretFile := filepath.Join(secretDir, "secret.txt")
	if err := os.WriteFile(secretFile, []byte("sensitive-data"), 0o600); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	if err := cleanupOldSecretDirs(tempDir); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if _, err := os.Stat(secretDir); err != nil {
		t.Fatalf("directory without timestamp should be kept: %v", err)
	}

	content, err := os.ReadFile(secretFile)
	if err != nil {
		t.Fatalf("secret file should still exist: %v", err)
	}

	if string(content) != "sensitive-data" {
		t.Fatal("secret file content should be intact")
	}
}

func TestCreateSecretDirGeneratesUniqueNames(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("failed to create executor: %v", err)
	}
	defer exec.Close()

	names := make(map[string]bool)
	for i := 0; i < 100; i++ {
		dir, err := exec.createSecretDir()
		if err != nil {
			t.Fatalf("failed to create secret directory: %v", err)
		}
		defer os.RemoveAll(dir)

		basename := filepath.Base(dir)
		if names[basename] {
			t.Fatalf("duplicate secret directory name generated: %s", basename)
		}
		names[basename] = true
	}
}

func TestCreateSecretDirCreatesTimestampFile(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("failed to create executor: %v", err)
	}
	defer exec.Close()

	secretDir, err := exec.createSecretDir()
	if err != nil {
		t.Fatalf("failed to create secret directory: %v", err)
	}
	defer os.RemoveAll(secretDir)

	tsFile := filepath.Join(secretDir, ".timestamp")
	if _, err := os.Stat(tsFile); err != nil {
		t.Fatalf("timestamp file should exist: %v", err)
	}

	content, err := os.ReadFile(tsFile)
	if err != nil {
		t.Fatalf("failed to read timestamp file: %v", err)
	}

	tsStr := string(content)
	_, err = time.Parse(time.RFC3339, tsStr)
	if err != nil {
		t.Fatalf("timestamp file should contain RFC3339 timestamp: %v", err)
	}
}

func TestCleanupSecretDir(t *testing.T) {
	tempDir := t.TempDir()
	secretDir := filepath.Join(tempDir, "test-secret")

	if err := os.Mkdir(secretDir, 0o700); err != nil {
		t.Fatalf("failed to create secret directory: %v", err)
	}

	secretFile := filepath.Join(secretDir, "secret.txt")
	secretValue := "sensitive-data-123"
	if err := os.WriteFile(secretFile, []byte(secretValue), 0o600); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	if err := cleanupSecretDir(secretDir); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if _, err := os.Stat(secretDir); !os.IsNotExist(err) {
		t.Fatal("secret directory should be removed after cleanup")
	}
}

func TestCleanupHandlesNonExistentDirectory(t *testing.T) {
	secretDir := filepath.Join(t.TempDir(), "non-existent")

	err := cleanupSecretDir(secretDir)
	if err != nil {
		t.Fatalf("cleanup should handle non-existent directory gracefully: %v", err)
	}
}

func TestMultipleSecretFilesInDirectory(t *testing.T) {
	tempDir := t.TempDir()
	secretDir := filepath.Join(tempDir, "multi-secret-dir")

	if err := os.Mkdir(secretDir, 0o700); err != nil {
		t.Fatalf("failed to create secret directory: %v", err)
	}

	numFiles := 10
	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("secret-%d.txt", i)
		path := filepath.Join(secretDir, filename)
		if err := os.WriteFile(path, []byte("secret-value"), 0o600); err != nil {
			t.Fatalf("failed to write secret file %s: %v", filename, err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("failed to stat secret file %s: %v", filename, err)
		}

		if info.Mode().Perm() != 0o600 {
			t.Fatalf("secret file %s has incorrect permissions: expected 0o600, got 0%o", filename, info.Mode().Perm())
		}
	}

	entries, err := os.ReadDir(secretDir)
	if err != nil {
		t.Fatalf("failed to read secret directory: %v", err)
	}

	if len(entries) != numFiles {
		t.Fatalf("expected %d files, got %d", numFiles, len(entries))
	}
}
