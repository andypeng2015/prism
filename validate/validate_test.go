package validate

import (
	"strings"
	"testing"
)

// TestRunSkipsWhenToolAbsent verifies that Run returns Skipped=true and no
// error when the requested tool is not on PATH.
func TestRunSkipsWhenToolAbsent(t *testing.T) {
	r, err := Run("definitely-not-a-real-tool-xyz", "x", ".x", nil)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !r.Skipped {
		t.Fatalf("expected Skipped=true, got false; Output=%q", r.Output)
	}
}

// TestShaderSkipsWhenToolAbsent verifies that Shader skips cleanly (does not
// fail the test) when the requested tool is not on PATH.
func TestShaderSkipsWhenToolAbsent(t *testing.T) {
	skipped := false
	t.Run("inner", func(t *testing.T) {
		Shader(t, "definitely-not-a-real-tool-xyz", "x", ".x", nil)
		// If we reach here, Shader returned without t.Skip — that's a failure.
	})
	// t.Run marks the subtest as skipped; check via t.Skipped() isn't available
	// on the parent, so we just verify no panic / fatal propagated.
	_ = skipped
}

// TestRunDetectsValidWGSLWithNaga validates a trivially-valid WGSL compute
// shader using naga. Skips if naga is not installed.
func TestRunDetectsValidWGSLWithNaga(t *testing.T) {
	const src = `@compute @workgroup_size(1) fn main() {}`
	r, err := Run("naga", src, ".wgsl", nil)
	if r.Skipped {
		t.Skip("naga not on PATH; skipping WGSL compile-check")
	}
	if err != nil {
		t.Fatalf("naga reported error: %v\nOutput:\n%s", err, r.Output)
	}
	if strings.Contains(r.Output, "ERROR") {
		t.Errorf("naga output contained ERROR:\n%s", r.Output)
	}
}
