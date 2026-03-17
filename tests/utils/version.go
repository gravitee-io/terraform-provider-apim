package utils

import (
	"os"
	"strings"
	"testing"
)

type ApimVersion int

const (
	ApimV4_9 ApimVersion = iota
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
	}
	return ""
}

func SkipFor(t *testing.T, version ...ApimVersion) {
	imageTag := os.Getenv("APIM_IMAGE_TAG")
	for _, v := range version {
		if strings.HasPrefix(imageTag, v.String()) {
			t.Skip("Skipping test for image tag" + imageTag + " as it does not support this feature")
			return
		}
	}
}
