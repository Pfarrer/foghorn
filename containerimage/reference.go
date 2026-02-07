package containerimage

import (
	"fmt"
	"strings"
)

type Reference struct {
	Repository string
	Tag        string
	Selector   Selector
}

func ParseReference(image string) (Reference, error) {
	if image == "" {
		return Reference{}, fmt.Errorf("image is required")
	}
	if strings.Contains(image, "@") {
		return Reference{}, fmt.Errorf("image digests are not supported")
	}

	repo, tag, err := splitTag(image)
	if err != nil {
		return Reference{}, err
	}
	if tag == "latest" {
		return Reference{}, fmt.Errorf("latest tag is not allowed")
	}
	selector, err := ParseSelector(tag)
	if err != nil {
		return Reference{}, err
	}
	return Reference{Repository: repo, Tag: tag, Selector: selector}, nil
}

func splitTag(image string) (string, string, error) {
	lastSlash := strings.LastIndex(image, "/")
	lastColon := strings.LastIndex(image, ":")
	if lastColon == -1 || lastColon < lastSlash+1 {
		return "", "", fmt.Errorf("image must include a version tag")
	}
	repo := image[:lastColon]
	tag := image[lastColon+1:]
	if repo == "" || tag == "" {
		return "", "", fmt.Errorf("image must include a version tag")
	}
	return repo, tag, nil
}
