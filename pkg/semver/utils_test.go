package semver

import (
	"testing"
)

func Test_bumpExistingNumeric(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		want     string
	}{
		{"empty string", "", "1"},
		{"no numeric suffix", "alpha", "alpha1"},
		{"simple numeric", "1", "2"},
		{"alpha with numeric", "alpha1", "alpha2"},
		{"complex with numeric", "alpha.beta.1", "alpha.beta.2"},
		{"larger number", "alpha99", "alpha100"},
		{"zero", "alpha0", "alpha1"},
		{"multiple numbers", "alpha1beta2", "alpha1beta3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bumpExistingNumeric(tt.existing); got != tt.want {
				t.Errorf("bumpExistingNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPrereleaseParts(t *testing.T) {
	tests := []struct {
		name        string
		prerelease  string
		wantPrefix  string
		wantNumeric string
	}{
		{"empty string", "", "", ""},
		{"no numeric part", "alpha", "alpha", ""},
		{"simple numeric", "1", "", "1"},
		{"alpha with numeric", "alpha1", "alpha", "1"},
		{"complex with numeric", "alpha.beta.1", "alpha.beta.", "1"},
		{"larger number", "alpha99", "alpha", "99"},
		{"zero", "alpha0", "alpha", "0"},
		{"multiple numbers", "alpha1beta2", "alpha1beta", "2"},
		{"dots and hyphens", "alpha.1-beta.2", "alpha.1-beta.", "2"},
		{"only dots", "...", "...", ""},
		{"only hyphens", "---", "---", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, gotNumeric := extractPrereleaseParts(tt.prerelease)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("extractPrereleaseParts() prefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotNumeric != tt.wantNumeric {
				t.Errorf("extractPrereleaseParts() numeric = %v, want %v", gotNumeric, tt.wantNumeric)
			}
		})
	}
}
