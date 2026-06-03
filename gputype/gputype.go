// Package gputype is the canonical GPU type vocabulary shared by all M31
// GPU-DSL compilers (Elio, Selena, Eos). Compilers map their own IR leaves
// onto these types; [m31labs.dev/prism/dialect] then spells each type in the
// target shading language (WGSL / GLSL / MSL).
//
// No external dependencies.
package gputype

import "fmt"

// Type is a shader value type. Concrete kinds implement it.
type Type interface{ isType() }

// Scalar is f32 / i32 / u32 / bool.
type Scalar struct{ Name string }

// Vec is an N-component vector of Elem (e.g. vec3<f32>).
type Vec struct {
	N    int
	Elem Scalar
}

// Mat is a matrix (e.g. mat4x4<f32>).
type Mat struct {
	Cols, Rows int
	Elem       Scalar
}

// Array is array<Elem, Len>; Len == 0 means a runtime-sized array.
type Array struct {
	Elem Type
	Len  int
}

func (Scalar) isType() {}
func (Vec) isType()    {}
func (Mat) isType()    {}
func (Array) isType()  {}

// Canonical scalars.
var (
	F32  = Scalar{"f32"}
	I32  = Scalar{"i32"}
	U32  = Scalar{"u32"}
	Bool = Scalar{"bool"}
)

// String returns a debug representation of the scalar (its name).
// This is NOT the backend spelling — use prism/dialect for that.
func (s Scalar) String() string { return s.Name }

// String returns a debug representation of the vector (e.g. "vec3<f32>").
// This is NOT the backend spelling — use prism/dialect for that.
func (v Vec) String() string { return fmt.Sprintf("vec%d<%s>", v.N, v.Elem.Name) }

// String returns a debug representation of the matrix (e.g. "mat4x4<f32>").
// This is NOT the backend spelling — use prism/dialect for that.
func (m Mat) String() string { return fmt.Sprintf("mat%dx%d<%s>", m.Cols, m.Rows, m.Elem.Name) }

// String returns a debug representation of the array (e.g. "[]f32" or "[6]f32").
// This is NOT the backend spelling — use prism/dialect for that.
func (a Array) String() string {
	if a.Len == 0 {
		return fmt.Sprintf("[]%s", a.Elem)
	}
	return fmt.Sprintf("[%d]%s", a.Len, a.Elem)
}
