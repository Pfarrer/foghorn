package executor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/logger"
)

func TestLoggerRedactsSecretsDuringCheckExecution(t *testing.T) {
	original := logger.GetGlobal()
	defer logger.SetGlobal(original)

	var buf bytes.Buffer
	l := logger.New(logger.LevelDebug, false)
	logger.SetGlobal(l)
	logger.SetOutput(&buf)

	secretValue := "super-secret-password-12345"
	logger.SetSecrets([]string{secretValue})

	logger.Info("Check test: Completed with env SMTP_PASSWORD=%s", secretValue)

	output := buf.String()
	if strings.Contains(output, secretValue) {
		t.Fatalf("secret should not appear in log output: %s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Fatalf("expected [REDACTED] in log output, got: %s", output)
	}
}

func TestLoggerRedactsAfterSecretsCleared(t *testing.T) {
	original := logger.GetGlobal()
	defer logger.SetGlobal(original)

	var buf bytes.Buffer
	l := logger.New(logger.LevelDebug, false)
	logger.SetGlobal(l)
	logger.SetOutput(&buf)

	secretValue := "clearable-secret-98765"
	logger.SetSecrets([]string{secretValue})
	logger.SetSecrets(nil)

	logger.Info("The secret is %s", secretValue)

	output := buf.String()
	if !strings.Contains(output, secretValue) {
		t.Fatalf("after clearing secrets, value should appear: %s", output)
	}
}

func TestValidateEnvNoSecrets(t *testing.T) {
	original := logger.GetGlobal()
	defer logger.SetGlobal(original)

	var buf bytes.Buffer
	l := logger.New(logger.LevelDebug, false)
	logger.SetGlobal(l)
	logger.SetOutput(&buf)

	exec := &DockerExecutor{}
	secrets := []string{"hunter2"}
	env := []string{
		"FOGHORN_CHECK_NAME=test",
		"PASSWORD=hunter2",
	}

	exec.validateEnvNoSecrets("test-check", env, secrets)

	output := buf.String()
	if strings.Contains(output, "hunter2") {
		t.Fatalf("warning should not contain raw secret: %s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Fatalf("warning should contain redacted value: %s", output)
	}
	if !strings.Contains(output, "[WARN]") {
		t.Fatalf("should produce a warning: %s", output)
	}
}

func TestValidateEnvNoSecretsClean(t *testing.T) {
	original := logger.GetGlobal()
	defer logger.SetGlobal(original)

	var buf bytes.Buffer
	l := logger.New(logger.LevelDebug, false)
	logger.SetGlobal(l)
	logger.SetOutput(&buf)

	exec := &DockerExecutor{}
	secrets := []string{"hunter2"}
	env := []string{
		"FOGHORN_CHECK_NAME=test",
		"PASSWORD_FILE=/run/foghorn/secrets/PASSWORD",
	}

	exec.validateEnvNoSecrets("test-check", env, secrets)

	output := buf.String()
	if strings.Contains(output, "[WARN]") {
		t.Fatalf("clean env should not produce warnings: %s", output)
	}
}

func TestBuildEnvVarsDoesNotExposeSecretsInEnv(t *testing.T) {
	exec := &DockerExecutor{
		secretResolver: &testSecretResolver{
			values: map[string]string{
				"secret://db/password": "db-secret-123",
			},
		},
	}

	checkConfig := &config.CheckConfig{
		Name:    "db-check",
		Image:   "test-image",
		Enabled: true,
		Env: map[string]string{
			"DB_PASSWORD": "secret://db/password",
			"DB_HOST":     "db.example.com",
		},
	}

	env, secretDir, _, err := exec.buildEnvVars(checkConfig)
	if err != nil {
		t.Fatalf("buildEnvVars failed: %v", err)
	}
	if secretDir != "" {
		defer cleanupSecretDir(secretDir)
	}

	for _, e := range env {
		if strings.Contains(e, "db-secret-123") {
			t.Fatalf("env should not contain raw secret: %s", e)
		}
	}
}

func TestEnvVarsSafeForProcessList(t *testing.T) {
	exec := &DockerExecutor{
		secretResolver: &testSecretResolver{
			values: map[string]string{
				"secret://smtp/pass": "smtp-password-xyz",
				"secret://api/token": "api-token-abc",
			},
		},
	}

	checkConfig := &config.CheckConfig{
		Name:    "mail-check",
		Image:   "test-image",
		Enabled: true,
		Env: map[string]string{
			"SMTP_PASSWORD": "secret://smtp/pass",
			"API_TOKEN":     "secret://api/token",
			"SMTP_HOST":     "smtp.example.com",
		},
	}

	env, secretDir, _, err := exec.buildEnvVars(checkConfig)
	if err != nil {
		t.Fatalf("buildEnvVars failed: %v", err)
	}
	if secretDir != "" {
		defer cleanupSecretDir(secretDir)
	}

	secrets := []string{"smtp-password-xyz", "api-token-abc"}
	for _, e := range env {
		for _, secret := range secrets {
			if strings.Contains(e, secret) {
				t.Fatalf("env var would expose secret in ps output: %s", e)
			}
		}
	}

	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	if envMap["SMTP_PASSWORD_FILE"] == "" {
		t.Fatal("expected SMTP_PASSWORD_FILE env var")
	}
	if !strings.HasPrefix(envMap["SMTP_PASSWORD_FILE"], "/run/foghorn/secrets/") {
		t.Fatalf("expected file path, got: %s", envMap["SMTP_PASSWORD_FILE"])
	}
	if envMap["API_TOKEN_FILE"] == "" {
		t.Fatal("expected API_TOKEN_FILE env var")
	}
}
