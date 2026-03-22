package sync

import "testing"

func TestIsValidSemver(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"1.19.2", true},
		{"1.21.0", true},
		{"1.20.4", true},
		{"1.21-pre1", false},
		{"1.20-rc2", false},
		{"1.20.4-3d", false},
		{"1.20.1-sodium", false},
		{"draft", false},
		{"3d", false},
		{"", false},
		{"1.19", false},
		{"1.19.2.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := IsValidSemver(tt.version); got != tt.want {
				t.Errorf("IsValidSemver(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
