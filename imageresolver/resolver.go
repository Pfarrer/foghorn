package imageresolver

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/image"
	"github.com/pfarrer/foghorn/containerimage"
)

type ImageLister interface {
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
}

type TagLister interface {
	ListTags(ctx context.Context, repository string) ([]string, error)
}

func Resolve(ctx context.Context, lister ImageLister, image string) (string, error) {
	return resolveWithTagLister(ctx, lister, image, newRegistryTagLister())
}

func resolveWithTagLister(ctx context.Context, _ ImageLister, image string, tags TagLister) (string, error) {
	ref, err := containerimage.ParseReference(image)
	if err != nil {
		return "", err
	}

	if ref.Selector.Kind == containerimage.SelectorFull {
		return image, nil
	}

	versions, err := availableVersions(ctx, tags, ref.Repository)
	if err != nil {
		return "", err
	}

	resolved, ok := containerimage.ResolveSelector(ref.Selector, versions)
	if !ok {
		return "", fmt.Errorf("no registry versions match selector %q for %s", ref.Tag, ref.Repository)
	}

	return fmt.Sprintf("%s:%s", ref.Repository, resolved.String()), nil
}

func availableVersions(ctx context.Context, tags TagLister, repo string) ([]containerimage.Version, error) {
	allTags, err := tags.ListTags(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to list registry tags: %w", err)
	}

	versions := make([]containerimage.Version, 0)
	for _, tag := range allTags {
		version, err := containerimage.ParseVersion(tag)
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}

	return versions, nil
}
