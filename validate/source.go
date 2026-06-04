package validate

import (
	"fmt"
	"strings"
)

// Backend names a shader or GPU-kernel source family that Prism can inspect.
type Backend string

const (
	BackendCUDA     Backend = "cuda"
	BackendMetal    Backend = "metal"
	BackendVulkan   Backend = "vulkan"
	BackendDirectML Backend = "directml"
	BackendWebGPU   Backend = "webgpu"
)

// Source is a backend-specific kernel or shader source blob.
type Source struct {
	Name    string
	Backend Backend
	Entry   string
	Source  string
}

// Extension returns the conventional file extension for backend sources.
func Extension(backend Backend) string {
	switch backend {
	case BackendCUDA:
		return ".cu"
	case BackendMetal:
		return ".metal"
	case BackendVulkan:
		return ".comp"
	case BackendDirectML:
		return ".hlsl"
	case BackendWebGPU:
		return ".wgsl"
	default:
		return ".txt"
	}
}

// CheckSource validates backend-neutral source invariants before tool-specific
// compilers are invoked. It is intentionally syntax-light so compilers can use
// it for generated placeholders as well as fully-native kernels.
func CheckSource(src Source) error {
	if src.Backend == "" {
		return fmt.Errorf("source backend is required")
	}
	if src.Entry == "" {
		return sourceError(src, "entry is required")
	}
	if src.Source == "" {
		return sourceError(src, "source is required")
	}
	needle, ok := entryNeedle(src.Backend, src.Entry)
	if !ok {
		return sourceError(src, fmt.Sprintf("unsupported backend %q", src.Backend))
	}
	if !strings.Contains(src.Source, needle) {
		return sourceError(src, fmt.Sprintf("%s source does not define entry %q", src.Backend, src.Entry))
	}
	return nil
}

// RunSource validates src with the conventional external tool for its backend
// when Prism knows one. Unsupported or unavailable tools return Skipped=true.
// CheckSource is always run first.
func RunSource(src Source) (Result, error) {
	if err := CheckSource(src); err != nil {
		return Result{}, err
	}
	tool, argv, ok := validatorForBackend(src.Backend)
	if !ok {
		return Result{Skipped: true}, nil
	}
	return Run(tool, src.Source, Extension(src.Backend), argv)
}

func validatorForBackend(backend Backend) (string, func(string) []string, bool) {
	switch backend {
	case BackendWebGPU:
		return "naga", nil, true
	case BackendVulkan:
		return "glslangValidator", func(path string) []string {
			return []string{"-V", path}
		}, true
	default:
		return "", nil, false
	}
}

func entryNeedle(backend Backend, entry string) (string, bool) {
	switch backend {
	case BackendCUDA:
		return "__global__ void " + entry + "(", true
	case BackendMetal:
		return "kernel void " + entry + "(", true
	case BackendVulkan:
		return "void " + entry + "(", true
	case BackendDirectML:
		return "eos_directml_graph " + entry + "(", true
	case BackendWebGPU:
		return "fn " + entry + "(", true
	default:
		return "", false
	}
}

func sourceError(src Source, msg string) error {
	if src.Name == "" {
		return fmt.Errorf("validate source: %s", msg)
	}
	return fmt.Errorf("validate source %q: %s", src.Name, msg)
}
