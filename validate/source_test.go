package validate

import (
	"strings"
	"testing"
)

func TestCheckSourceAcceptsBackendEntryShapes(t *testing.T) {
	tests := []struct {
		name    string
		backend Backend
		entry   string
		source  string
	}{
		{"cuda", BackendCUDA, "step_cuda", `extern "C" __global__ void step_cuda(float *out) {}`},
		{"metal", BackendMetal, "step_metal", `kernel void step_metal(device float *out) {}`},
		{"vulkan", BackendVulkan, "step_vulkan", `#version 450
void step_vulkan() {}`},
		{"directml", BackendDirectML, "step_directml", `eos_directml_graph step_directml() {}`},
		{"legacy directml", BackendDirectML, "step_directml", `manta_directml_graph step_directml() {}`},
		{"webgpu", BackendWebGPU, "step_webgpu", `@compute @workgroup_size(1)
fn step_webgpu() {}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckSource(Source{Name: "k", Backend: tt.backend, Entry: tt.entry, Source: tt.source}); err != nil {
				t.Fatalf("CheckSource() error = %v", err)
			}
		})
	}
}

func TestCheckSourceRejectsMissingEntry(t *testing.T) {
	err := CheckSource(Source{
		Name:    "score",
		Backend: BackendWebGPU,
		Entry:   "score_webgpu",
		Source:  `@compute @workgroup_size(1) fn other() {}`,
	})
	if err == nil {
		t.Fatal("CheckSource() succeeded, want missing entry error")
	}
	if !strings.Contains(err.Error(), `source "score"`) || !strings.Contains(err.Error(), `entry "score_webgpu"`) {
		t.Fatalf("CheckSource() error = %v", err)
	}
}

func TestRunSourceSkipsWhenNoExternalToolIsKnown(t *testing.T) {
	res, err := RunSource(Source{
		Name:    "score",
		Backend: BackendCUDA,
		Entry:   "score_cuda",
		Source:  `extern "C" __global__ void score_cuda(float *out) {}`,
	})
	if err != nil {
		t.Fatalf("RunSource() error = %v", err)
	}
	if !res.Skipped {
		t.Fatalf("RunSource() Skipped = false, want true for CUDA without configured external tool")
	}
}
