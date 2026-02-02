package main

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
)

func TestVerifyImageAvailability_NoChecks(t *testing.T) {
	cfg := &config.Config{
		Checks: []config.CheckConfig{},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err != nil {
		t.Errorf("Expected no error with no checks, got: %v", err)
	}
}

func TestVerifyImageAvailability_NoEnabledChecks(t *testing.T) {
	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "disabled-check",
				Image:   "test/image:latest",
				Enabled: false,
			},
		},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err != nil {
		t.Errorf("Expected no error with no enabled checks, got: %v", err)
	}
}

func TestVerifyImageAvailability_MissingImage(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "test-check",
				Image:   "this-image-definitely-does-not-exist-12345:latest",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	expected := "Error: The following Docker images are not available locally"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error message to contain '%s', got: %v", expected, err)
	}

	expectedPull := "docker pull this-image-definitely-does-not-exist-12345:latest"
	if err != nil && !containsString(err.Error(), expectedPull) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedPull, err)
	}
}

func TestVerifyImageAvailability_ExistingImage(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	_, _, err = cli.ImageInspectWithRaw(context.Background(), "alpine:latest")
	if err != nil {
		if client.IsErrNotFound(err) {
			t.Skipf("Skipping test: alpine:latest not available locally. Run 'docker pull alpine:latest'")
		}
		t.Skipf("Skipping test: error checking alpine:latest: %v", err)
	}

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "alpine-check",
				Image:   "alpine:latest",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err != nil {
		t.Errorf("Expected no error for existing image, got: %v", err)
	}
}

func TestVerifyImageAvailability_MultipleChecksSameImage(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "check1",
				Image:   "this-image-definitely-does-not-exist-12345:latest",
				Enabled: true,
			},
			{
				Name:    "check2",
				Image:   "this-image-definitely-does-not-exist-12345:latest",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	if err != nil && !containsString(err.Error(), "(required by: check1, check2)") &&
		!containsString(err.Error(), "(required by: check2, check1)") {
		t.Errorf("Expected error to list both check names, got: %v", err)
	}
}

func TestVerifyImageAvailability_MultipleImages(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "check1",
				Image:   "this-image-definitely-does-not-exist-12345:latest",
				Enabled: true,
			},
			{
				Name:    "check2",
				Image:   "another-missing-image-67890:latest",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing images, got nil")
	}

	expected := "this-image-definitely-does-not-exist-12345:latest"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got: %v", expected, err)
	}

	expected2 := "another-missing-image-67890:latest"
	if err != nil && !containsString(err.Error(), expected2) {
		t.Errorf("Expected error to contain '%s', got: %v", expected2, err)
	}
}

func TestVerifyImageAvailability_MixedMissingAndExisting(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	_, _, err = cli.ImageInspectWithRaw(context.Background(), "alpine:latest")
	if err != nil {
		if client.IsErrNotFound(err) {
			t.Skipf("Skipping test: alpine:latest not available locally. Run 'docker pull alpine:latest'")
		}
		t.Skipf("Skipping test: error checking alpine:latest: %v", err)
	}

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "existing-check",
				Image:   "alpine:latest",
				Enabled: true,
			},
			{
				Name:    "missing-check",
				Image:   "this-image-definitely-does-not-exist-12345:latest",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	expected := "this-image-definitely-does-not-exist-12345:latest"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got: %v", expected, err)
	}

	if err != nil && containsString(err.Error(), "alpine:latest") {
		t.Errorf("Error should not mention existing image alpine:latest, got: %v", err)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
