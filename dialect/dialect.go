// Package dialect provides the per-backend spelling layer for GPU shader
// codegen. It is IR-agnostic: it knows only [gputype.Type] and produces the
// correct surface string for each target shading language.
//
// Lifted from Selena's per-backend dialects so that Selena (Task 3) and Elio
// (Task 4) can delete their copies and import this package instead.
package dialect

import "m31labs.dev/prism/gputype"

// Dialect spells the backend-specific, IR-agnostic parts of shader codegen.
type Dialect interface {
	TypeName(gputype.Type) string              // WGSL vec3<f32> · GLSL vec3 · Metal float3
	Builtin(name string, args []string) string // a builtin call, e.g. dot(a, b) (Metal may rename)
	Swizzle(expr, components string) string    // expr.xyz (same on all backends today)
	Sample(tex, uv string) string              // WGSL textureSample(...) · GLSL texture2D(...) · Metal t.sample(...)
	// SampleCube renders a cube-map texture sample by a vec3 direction vector.
	// The cube texture must be declared with the appropriate cube-map type for the
	// target backend (texture_cube<f32>, samplerCube, or texturecube<float>).
	// WGSL:         textureSample(tex, texSampler, dir)
	// GLSL ES 1.00: textureCube(tex, dir)
	// GLSL ES 3.00: texture(tex, dir)
	// Metal:        tex.sample(texSampler, dir)
	SampleCube(tex, dir string) string
	// Ternary renders a conditional expression. cond is bool, then and alt have
	// the same type. WGSL has no ? : operator — it emits select(alt, then, cond).
	// GLSL, GLES, and Metal use the standard (cond ? then : alt) form.
	Ternary(cond, then, alt string) string
}
