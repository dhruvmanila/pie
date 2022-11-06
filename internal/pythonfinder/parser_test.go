package pythonfinder

import (
	"errors"
	"reflect"
	"testing"
)

func TestInvalidVersion(t *testing.T) {
	for _, version := range []string{"", "3", "3.9", "3.10", "3.12.0a1"} {
		t.Run(version, func(t *testing.T) {
			got, err := parseVersion(version)
			if !errors.Is(err, ErrInvalidVersion) {
				t.Errorf("parseVersion(%q) = %v, want %v", version, got, ErrInvalidVersion)
			}
		})
	}
}

func TestValidVersion(t *testing.T) {
	tests := map[string]*VersionInfo{
		"3.9.0":  {Major: 3, Minor: 9, Patch: 0},
		"3.10.5": {Major: 3, Minor: 10, Patch: 5},
	}

	for version, want := range tests {
		t.Run(version, func(t *testing.T) {
			got, err := parseVersion(version)
			if err != nil {
				t.Errorf("parseVersion(%q) error = %v", version, err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("parseVersion(%q) = %v, want %v", version, got, want)
			}
		})
	}
}
