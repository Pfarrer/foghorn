package containerimage

import "testing"

func TestParseSelector(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		kind    SelectorKind
		wantErr bool
	}{
		{name: "major", tag: "1", kind: SelectorMajor},
		{name: "major patch", tag: "2.5", kind: SelectorMajorPatch},
		{name: "full", tag: "3.4.5", kind: SelectorFull},
		{name: "invalid", tag: "1.2.3.4", wantErr: true},
		{name: "invalid chars", tag: "1.x", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, err := ParseSelector(tt.tag)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if selector.Kind != tt.kind {
				t.Fatalf("expected kind %v, got %v", tt.kind, selector.Kind)
			}
		})
	}
}

func TestResolveSelector(t *testing.T) {
	versions := []Version{
		{Major: 1, Minor: 0, Patch: 1},
		{Major: 1, Minor: 2, Patch: 1},
		{Major: 1, Minor: 1, Patch: 2},
		{Major: 2, Minor: 0, Patch: 0},
	}

	selector, err := ParseSelector("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	best, ok := ResolveSelector(selector, versions)
	if !ok {
		t.Fatalf("expected match")
	}
	if best.String() != "1.2.1" {
		t.Fatalf("expected 1.2.1, got %s", best.String())
	}

	selector, err = ParseSelector("1.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	best, ok = ResolveSelector(selector, versions)
	if !ok {
		t.Fatalf("expected match")
	}
	if best.String() != "1.1.2" {
		t.Fatalf("expected 1.1.2, got %s", best.String())
	}
}
