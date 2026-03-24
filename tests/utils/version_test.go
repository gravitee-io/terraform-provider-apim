package utils

import (
	"os"
	"testing"
)

func TestApimVersion_String(t *testing.T) {
	tests := []struct {
		v    ApimVersion
		want string
	}{
		{ApimV4_9, "4.9"},
		{ApimV4_10, "4.10"},
		{ApimV4_11, "4.11"},
		{ApimUnknown, "unknown"},
		{ApimVersion(99), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("ApimVersion.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseApimVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    ApimVersion
	}{
		{
			name:    "parse 4.9",
			version: "4.9",
			want:    ApimV4_9,
		},
		{
			name:    "parse 4.10",
			version: "4.10",
			want:    ApimV4_10,
		},
		{
			name:    "parse 4.11",
			version: "4.11",
			want:    ApimV4_11,
		},
		{
			name:    "parse 4.9 with patch",
			version: "4.9.1",
			want:    ApimV4_9,
		},
		{
			name:    "parse 4.10 with rc",
			version: "4.10.0-rc1",
			want:    ApimV4_10,
		},
		{
			name:    "parse unknown version",
			version: "4.12.0",
			want:    ApimUnknown,
		},
		{
			name:    "parse empty string",
			version: "",
			want:    ApimUnknown,
		},
		{
			name:    "parse invalid version",
			version: "invalid",
			want:    ApimUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseApimVersion(tt.version); got != tt.want {
				t.Errorf("ParseApimVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipFor(t *testing.T) {
	origEnv := os.Getenv("APIM_IMAGE_TAG")
	defer os.Setenv("APIM_IMAGE_TAG", origEnv)

	tests := []struct {
		name        string
		imageTag    string
		skipFor     []ApimVersion
		wantSkipped bool
	}{
		{
			name:        "matches version 4.9",
			imageTag:    "4.9.1",
			skipFor:     []ApimVersion{ApimV4_9},
			wantSkipped: true,
		},
		{
			name:        "matches version 4.10",
			imageTag:    "4.10.0-rc1",
			skipFor:     []ApimVersion{ApimV4_10},
			wantSkipped: true,
		},
		{
			name:        "no match",
			imageTag:    "4.12.0",
			skipFor:     []ApimVersion{ApimV4_9, ApimV4_10, ApimV4_11},
			wantSkipped: false,
		},
		{
			name:        "empty tag",
			imageTag:    "",
			skipFor:     []ApimVersion{ApimV4_9},
			wantSkipped: false,
		},
		{
			name:        "multiple skip for, matches one",
			imageTag:    "4.11.2",
			skipFor:     []ApimVersion{ApimV4_9, ApimV4_11},
			wantSkipped: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipFor(tt.imageTag, tt.skipFor...)
			if got != tt.wantSkipped {
				t.Errorf("shouldSkipFor() = %v, want %v", got, tt.wantSkipped)
			}
		})
	}
}
