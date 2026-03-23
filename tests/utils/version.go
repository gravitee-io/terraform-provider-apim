package utils

import (
	"os"
	"strings"
	"testing"
)

type ApimVersion int

const (
	ApimUnknown ApimVersion = iota
	ApimV4_9
	ApimV4_10
	ApimV4_11
)

func (v ApimVersion) String() string {
	switch v {
	case ApimV4_9:
		return "4.9"
	case ApimV4_10:
		return "4.10"
	case ApimV4_11:
		return "4.11"
	case ApimUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

func ParseApimVersion(s string) ApimVersion {
	switch {
	case strings.HasSuffix(s, "4.9"):
		return ApimV4_9
	case strings.HasSuffix(s, "4.10"):
		return ApimV4_10
	case strings.HasSuffix(s, "4.11"):
		return ApimV4_11
	default:
		return ApimUnknown
	}
}

func SkipFor(t *testing.T, version ...ApimVersion) {
	imageTag := os.Getenv("APIM_IMAGE_TAG")
	for _, v := range version {
		if strings.HasPrefix(imageTag, v.String()) {
			t.Skip("Skipping test for image tag" + imageTag + " as it contains unsupported content")
			return
		}
	}
}
