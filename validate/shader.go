package validate

import (
	"strings"
	"testing"
)

// Shader runs Run and adapts the result to a test:
//   - Skipped=true  → t.Skip (tool not on PATH)
//   - non-nil error → t.Fatalf
//   - Output contains "ERROR" → t.Errorf
func Shader(t *testing.T, tool, source, ext string, argv func(path string) []string) {
	t.Helper()
	r, err := Run(tool, source, ext, argv)
	if r.Skipped {
		t.Skipf("%s: skipped (not on PATH)", tool)
		return
	}
	if err != nil {
		t.Fatalf("%s validation failed: %v\n%s", tool, err, r.Output)
	}
	if strings.Contains(r.Output, "ERROR") {
		t.Errorf("%s reported errors:\n%s", tool, r.Output)
	}
}
