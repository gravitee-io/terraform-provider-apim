package acceptance_test

import (
	"sync"
	"testing"
)

var portalAcceptanceTestMu sync.Mutex

func lockPortalAcceptanceTest(t *testing.T) {
	t.Helper()
	portalAcceptanceTestMu.Lock()
	t.Cleanup(portalAcceptanceTestMu.Unlock)
}
