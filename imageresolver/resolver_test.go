package imageresolver

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/image"
)

type stubLister struct {
	images []image.Summary
}

func (s stubLister) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return s.images, nil
}

func TestResolve(t *testing.T) {
	lister := stubLister{images: []image.Summary{
		{RepoTags: []string{"repo/check:1.0.0", "repo/check:1.2.0", "repo/check:1.1.3"}},
		{RepoTags: []string{"repo/other:2.0.0"}},
	}}

	resolved, err := Resolve(context.Background(), lister, "repo/check:1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "repo/check:1.2.0" {
		t.Fatalf("expected repo/check:1.2.0, got %s", resolved)
	}

	resolved, err = Resolve(context.Background(), lister, "repo/check:1.3.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "repo/check:1.3.0" {
		t.Fatalf("expected repo/check:1.3.0, got %s", resolved)
	}
}

func TestResolveMajorPatch(t *testing.T) {
	lister := stubLister{images: []image.Summary{
		{RepoTags: []string{"repo/check:1.0.2", "repo/check:1.4.2", "repo/check:1.1.1"}},
	}}

	resolved, err := Resolve(context.Background(), lister, "repo/check:1.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "repo/check:1.4.2" {
		t.Fatalf("expected repo/check:1.4.2, got %s", resolved)
	}
}
