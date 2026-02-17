package imageresolver

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types/image"
)

type stubLister struct {
	images []image.Summary
}

func (s stubLister) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return s.images, nil
}

type stubTagLister struct {
	tagsByRepository map[string][]string
	errByRepository  map[string]error
}

func (s stubTagLister) ListTags(ctx context.Context, repository string) ([]string, error) {
	if err, ok := s.errByRepository[repository]; ok {
		return nil, err
	}
	return s.tagsByRepository[repository], nil
}

func TestResolve(t *testing.T) {
	lister := stubLister{}
	tags := stubTagLister{
		tagsByRepository: map[string][]string{
			"repo/check": {"1.0.0", "1.2.0", "1.1.3", "bad-tag", "2.0.0"},
		},
	}

	resolved, err := resolveWithTagLister(context.Background(), lister, "repo/check:1", tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "repo/check:1.2.0" {
		t.Fatalf("expected repo/check:1.2.0, got %s", resolved)
	}

	resolved, err = resolveWithTagLister(context.Background(), lister, "repo/check:1.3.0", tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "repo/check:1.3.0" {
		t.Fatalf("expected repo/check:1.3.0, got %s", resolved)
	}
}

func TestResolveMajorPatch(t *testing.T) {
	lister := stubLister{}
	tags := stubTagLister{
		tagsByRepository: map[string][]string{
			"repo/check": {"1.0.2", "1.4.2", "1.1.1"},
		},
	}

	resolved, err := resolveWithTagLister(context.Background(), lister, "repo/check:1.2", tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "repo/check:1.4.2" {
		t.Fatalf("expected repo/check:1.4.2, got %s", resolved)
	}
}

func TestResolvePartialNoMatch(t *testing.T) {
	lister := stubLister{}
	tags := stubTagLister{
		tagsByRepository: map[string][]string{
			"repo/check": {"2.0.0", "2.1.0"},
		},
	}

	_, err := resolveWithTagLister(context.Background(), lister, "repo/check:1", tags)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != `no registry versions match selector "1" for repo/check` {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePartialListTagsError(t *testing.T) {
	lister := stubLister{}
	tags := stubTagLister{
		errByRepository: map[string]error{
			"repo/check": fmt.Errorf("boom"),
		},
	}

	_, err := resolveWithTagLister(context.Background(), lister, "repo/check:1", tags)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to list registry tags: boom" {
		t.Fatalf("unexpected error: %v", err)
	}
}
