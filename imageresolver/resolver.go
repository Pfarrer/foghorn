package imageresolver

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/pfarrer/foghorn/containerimage"
)

type ImageLister interface {
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
}

func Resolve(ctx context.Context, lister ImageLister, image string) (string, error) {
	ref, err := containerimage.ParseReference(image)
	if err != nil {
		return "", err
	}

	if ref.Selector.Kind == containerimage.SelectorFull {
		return image, nil
	}

	versions, err := availableVersions(ctx, lister, ref.Repository)
	if err != nil {
		return "", err
	}

	resolved, ok := containerimage.ResolveSelector(ref.Selector, versions)
	if !ok {
		return "", fmt.Errorf("no local versions match selector %q for %s", ref.Tag, ref.Repository)
	}

	return fmt.Sprintf("%s:%s", ref.Repository, resolved.String()), nil
}

func availableVersions(ctx context.Context, lister ImageLister, repo string) ([]containerimage.Version, error) {
	images, err := lister.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list local images: %w", err)
	}

	versions := make([]containerimage.Version, 0)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			repoTag := tag
			if !strings.HasPrefix(repoTag, repo+":") {
				continue
			}
			versionTag := strings.TrimPrefix(repoTag, repo+":")
			version, err := containerimage.ParseVersion(versionTag)
			if err != nil {
				continue
			}
			versions = append(versions, version)
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) > 0
	})

	return versions, nil
}
