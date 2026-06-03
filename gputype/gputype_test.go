package gputype_test

import (
	"testing"

	"m31labs.dev/prism/gputype"
)

// Compile-time interface checks.
var (
	_ gputype.Type = gputype.Scalar{}
	_ gputype.Type = gputype.Vec{}
	_ gputype.Type = gputype.Mat{}
	_ gputype.Type = gputype.Array{}
)

func TestScalarString(t *testing.T) {
	if got := gputype.F32.String(); got != "f32" {
		t.Errorf("F32.String() = %q, want %q", got, "f32")
	}
	if got := gputype.I32.String(); got != "i32" {
		t.Errorf("I32.String() = %q, want %q", got, "i32")
	}
	if got := gputype.U32.String(); got != "u32" {
		t.Errorf("U32.String() = %q, want %q", got, "u32")
	}
	if got := gputype.Bool.String(); got != "bool" {
		t.Errorf("Bool.String() = %q, want %q", got, "bool")
	}
}

func TestVecString(t *testing.T) {
	v := gputype.Vec{N: 3, Elem: gputype.F32}
	if got := v.String(); got != "vec3<f32>" {
		t.Errorf("Vec{3,F32}.String() = %q, want %q", got, "vec3<f32>")
	}
}

func TestMatString(t *testing.T) {
	m := gputype.Mat{Cols: 4, Rows: 4, Elem: gputype.F32}
	if got := m.String(); got != "mat4x4<f32>" {
		t.Errorf("Mat{4,4,F32}.String() = %q, want %q", got, "mat4x4<f32>")
	}
}

func TestArrayString(t *testing.T) {
	runtimeSized := gputype.Array{Elem: gputype.F32, Len: 0}
	if got := runtimeSized.String(); got != "[]f32" {
		t.Errorf("Array{F32,0}.String() = %q, want %q", got, "[]f32")
	}

	fixed := gputype.Array{Elem: gputype.F32, Len: 6}
	if got := fixed.String(); got != "[6]f32" {
		t.Errorf("Array{F32,6}.String() = %q, want %q", got, "[6]f32")
	}
}
