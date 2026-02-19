package daemon

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/containerimage"
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
				Image:   "test/image:1.0.0",
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
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
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

	expectedPull := "docker pull this-image-definitely-does-not-exist-12345:1.2.3"
	if err != nil && !containsString(err.Error(), expectedPull) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedPull, err)
	}
}

func TestVerifyImageAvailability_ExistingImage(t *testing.T) {
	image := findLocalSemverImage(t)

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "alpine-check",
				Image:   image,
				Enabled: true,
			},
		},
	}

	err := verifyImageAvailabilityFn(cfg)
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
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
			{
				Name:    "check2",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
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
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
			{
				Name:    "check2",
				Image:   "another-missing-image-67890:2.3.4",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing images, got nil")
	}

	expected := "this-image-definitely-does-not-exist-12345:1.2.3"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got: %v", expected, err)
	}

	expected2 := "another-missing-image-67890:2.3.4"
	if err != nil && !containsString(err.Error(), expected2) {
		t.Errorf("Expected error to contain '%s', got: %v", expected2, err)
	}
}

func TestVerifyImageAvailability_MixedMissingAndExisting(t *testing.T) {
	image := findLocalSemverImage(t)

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "existing-check",
				Image:   image,
				Enabled: true,
			},
			{
				Name:    "missing-check",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
		},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	expected := "this-image-definitely-does-not-exist-12345:1.2.3"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got: %v", expected, err)
	}

	if err != nil && containsString(err.Error(), image) {
		t.Errorf("Error should not mention existing image %s, got: %v", image, err)
	}
}

func findLocalSemverImage(t *testing.T) string {
	t.Helper()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		t.Skipf("Skipping test: cannot list docker images: %v", err)
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			ref, err := containerimage.ParseReference(tag)
			if err != nil {
				continue
			}
			if ref.Selector.Kind != containerimage.SelectorFull {
				continue
			}
			return tag
		}
	}

	t.Skipf("Skipping test: no local images with semantic version tags found")
	return ""
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
