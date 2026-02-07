package containerimage

import (
	"fmt"
	"strconv"
	"strings"
)

type SelectorKind int

const (
	SelectorMajor SelectorKind = iota
	SelectorMajorPatch
	SelectorFull
)

type Selector struct {
	Kind  SelectorKind
	Major int
	Minor int
	Patch int
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func ParseSelector(tag string) (Selector, error) {
	if tag == "" {
		return Selector{}, fmt.Errorf("tag is required")
	}
	parts := strings.Split(tag, ".")
	switch len(parts) {
	case 1:
		major, err := parsePart(parts[0])
		if err != nil {
			return Selector{}, err
		}
		return Selector{Kind: SelectorMajor, Major: major}, nil
	case 2:
		major, err := parsePart(parts[0])
		if err != nil {
			return Selector{}, err
		}
		patch, err := parsePart(parts[1])
		if err != nil {
			return Selector{}, err
		}
		return Selector{Kind: SelectorMajorPatch, Major: major, Patch: patch}, nil
	case 3:
		major, err := parsePart(parts[0])
		if err != nil {
			return Selector{}, err
		}
		minor, err := parsePart(parts[1])
		if err != nil {
			return Selector{}, err
		}
		patch, err := parsePart(parts[2])
		if err != nil {
			return Selector{}, err
		}
		return Selector{Kind: SelectorFull, Major: major, Minor: minor, Patch: patch}, nil
	default:
		return Selector{}, fmt.Errorf("invalid tag format")
	}
}

func ParseVersion(tag string) (Version, error) {
	parts := strings.Split(tag, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid version format")
	}
	major, err := parsePart(parts[0])
	if err != nil {
		return Version{}, err
	}
	minor, err := parsePart(parts[1])
	if err != nil {
		return Version{}, err
	}
	patch, err := parsePart(parts[2])
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}
	return 0
}

func (s Selector) Matches(v Version) bool {
	switch s.Kind {
	case SelectorMajor:
		return v.Major == s.Major
	case SelectorMajorPatch:
		return v.Major == s.Major && v.Patch == s.Patch
	case SelectorFull:
		return v.Major == s.Major && v.Minor == s.Minor && v.Patch == s.Patch
	default:
		return false
	}
}

func ResolveSelector(selector Selector, versions []Version) (Version, bool) {
	var best *Version
	for _, v := range versions {
		if !selector.Matches(v) {
			continue
		}
		if best == nil || v.Compare(*best) > 0 {
			copy := v
			best = &copy
		}
	}
	if best == nil {
		return Version{}, false
	}
	return *best, true
}

func parsePart(part string) (int, error) {
	if part == "" {
		return 0, fmt.Errorf("invalid tag format")
	}
	for _, r := range part {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("invalid tag format")
		}
	}
	value, err := strconv.Atoi(part)
	if err != nil {
		return 0, fmt.Errorf("invalid tag format")
	}
	return value, nil
}
