// Package validate provides a thin harness for running offline shader-compiler
// tools (naga, glslangValidator, …) as subprocesses. It is intentionally
// testing-free so it can be used from both test and non-test code; the
// testing-aware helper lives in shader.go.
package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Result is the outcome of running a validator.
type Result struct {
	Skipped bool   // tool not on PATH — validation was skipped, not failed
	Output  string // combined stdout+stderr from the tool
}

// Run writes source to a temp file named "shader"+ext, then runs tool with
// argv(path) as arguments (cmd.Dir set to the temp dir so stray outputs are
// contained). If tool is not on PATH it returns Result{Skipped:true}, nil.
// If argv is nil, the file path is the single argument.
// A non-nil error means the tool exited non-zero; Output is always populated.
func Run(tool, source, ext string, argv func(path string) []string) (Result, error) {
	bin, err := exec.LookPath(tool)
	if err != nil {
		return Result{Skipped: true}, nil
	}

	dir, err := os.MkdirTemp("", "prism-validate-*")
	if err != nil {
		return Result{}, fmt.Errorf("validate: create temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "shader"+ext)
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		return Result{}, fmt.Errorf("validate: write shader: %w", err)
	}

	var args []string
	if argv != nil {
		args = argv(path)
	} else {
		args = []string{path}
	}

	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	res := Result{Output: string(out)}
	if err != nil {
		return res, fmt.Errorf("validate: %s: %w", tool, err)
	}
	return res, nil
}
