package dialect_test

import (
	"testing"

	"m31labs.dev/prism/dialect"
	"m31labs.dev/prism/gputype"
)

// Each row in the table drives one assertion against one dialect method.
type row struct {
	label string
	got   string
	want  string
}

func run(t *testing.T, rows []row) {
	t.Helper()
	for _, r := range rows {
		if r.got != r.want {
			t.Errorf("%s: got %q, want %q", r.label, r.got, r.want)
		}
	}
}

// ---------------------------------------------------------------------------
// WGSL
// ---------------------------------------------------------------------------

func TestWGSL(t *testing.T) {
	d := dialect.WGSL{}

	// TypeName — lifted from selena/emit/wgsl/wgsl.go:110-127
	// and elio/emit/wgsl/wgsl.go:199-217 for Array.
	run(t, []row{
		{"WGSL TypeName f32", d.TypeName(gputype.F32), "f32"},
		{"WGSL TypeName i32", d.TypeName(gputype.I32), "i32"},
		{"WGSL TypeName u32", d.TypeName(gputype.U32), "u32"},
		{"WGSL TypeName bool", d.TypeName(gputype.Bool), "bool"},
		{"WGSL TypeName vec2<f32>", d.TypeName(gputype.Vec{N: 2, Elem: gputype.F32}), "vec2<f32>"},
		{"WGSL TypeName vec3<f32>", d.TypeName(gputype.Vec{N: 3, Elem: gputype.F32}), "vec3<f32>"},
		{"WGSL TypeName vec4<f32>", d.TypeName(gputype.Vec{N: 4, Elem: gputype.F32}), "vec4<f32>"},
		{"WGSL TypeName vec4<i32>", d.TypeName(gputype.Vec{N: 4, Elem: gputype.I32}), "vec4<i32>"},
		{"WGSL TypeName vec4<u32>", d.TypeName(gputype.Vec{N: 4, Elem: gputype.U32}), "vec4<u32>"},
		{"WGSL TypeName mat3x3<f32>", d.TypeName(gputype.Mat{Cols: 3, Rows: 3, Elem: gputype.F32}), "mat3x3<f32>"},
		{"WGSL TypeName mat4x4<f32>", d.TypeName(gputype.Mat{Cols: 4, Rows: 4, Elem: gputype.F32}), "mat4x4<f32>"},
		{"WGSL TypeName mat4x3<f32>", d.TypeName(gputype.Mat{Cols: 4, Rows: 3, Elem: gputype.F32}), "mat4x3<f32>"},
		// Array — lifted from elio/emit/wgsl/wgsl.go:207-210
		{"WGSL TypeName array<f32>", d.TypeName(gputype.Array{Elem: gputype.F32, Len: 0}), "array<f32>"},
		{"WGSL TypeName array<f32,6>", d.TypeName(gputype.Array{Elem: gputype.F32, Len: 6}), "array<f32, 6>"},
		{"WGSL TypeName array<u32,4>", d.TypeName(gputype.Array{Elem: gputype.U32, Len: 4}), "array<u32, 4>"},
	})

	// Builtin — spell.Call with WGSL builtinSpellings (all identity);
	// lifted from selena/emit/internal/spell/spell.go:10-16 and
	// selena/emit/wgsl/wgsl.go:137-164.
	run(t, []row{
		{"WGSL Builtin dot", d.Builtin("dot", []string{"a", "b"}), "dot(a, b)"},
		{"WGSL Builtin mix", d.Builtin("mix", []string{"a", "b", "t"}), "mix(a, b, t)"},
		{"WGSL Builtin normalize", d.Builtin("normalize", []string{"v"}), "normalize(v)"},
		{"WGSL Builtin clamp", d.Builtin("clamp", []string{"x", "lo", "hi"}), "clamp(x, lo, hi)"},
		{"WGSL Builtin unknown", d.Builtin("myFunc", []string{"x"}), "myFunc(x)"},
	})

	// Swizzle — from selena/ir/ir.go:152 `Print(x.E, d) + "." + x.Field`
	run(t, []row{
		{"WGSL Swizzle xyz", d.Swizzle("v", "xyz"), "v.xyz"},
		{"WGSL Swizzle rgb", d.Swizzle("col", "rgb"), "col.rgb"},
	})

	// Sample — lifted from selena/emit/wgsl/wgsl.go:106-108
	// "textureSample(" + tex + ", " + tex + "Sampler, " + uv + ")"
	run(t, []row{
		{"WGSL Sample", d.Sample("t", "uv"), "textureSample(t, tSampler, uv)"},
		{"WGSL Sample albedo", d.Sample("albedo", "coord"), "textureSample(albedo, albedoSampler, coord)"},
	})
}

// ---------------------------------------------------------------------------
// GLSL
// ---------------------------------------------------------------------------

func TestGLSL(t *testing.T) {
	d := dialect.GLSL{}

	// TypeName — lifted from selena/emit/glsl/glsl.go:76-93
	// and elio/emit/glsl/glsl.go:442-480 for i32/u32 vectors, non-square mats.
	run(t, []row{
		{"GLSL TypeName float", d.TypeName(gputype.F32), "float"},
		{"GLSL TypeName int", d.TypeName(gputype.I32), "int"},
		{"GLSL TypeName uint", d.TypeName(gputype.U32), "uint"},
		{"GLSL TypeName vec2", d.TypeName(gputype.Vec{N: 2, Elem: gputype.F32}), "vec2"},
		{"GLSL TypeName vec3", d.TypeName(gputype.Vec{N: 3, Elem: gputype.F32}), "vec3"},
		{"GLSL TypeName vec4", d.TypeName(gputype.Vec{N: 4, Elem: gputype.F32}), "vec4"},
		{"GLSL TypeName ivec3", d.TypeName(gputype.Vec{N: 3, Elem: gputype.I32}), "ivec3"},
		{"GLSL TypeName uvec4", d.TypeName(gputype.Vec{N: 4, Elem: gputype.U32}), "uvec4"},
		{"GLSL TypeName mat3", d.TypeName(gputype.Mat{Cols: 3, Rows: 3, Elem: gputype.F32}), "mat3"},
		{"GLSL TypeName mat4", d.TypeName(gputype.Mat{Cols: 4, Rows: 4, Elem: gputype.F32}), "mat4"},
		// Non-square mat from elio/emit/glsl/glsl.go:467-468
		{"GLSL TypeName mat4x3", d.TypeName(gputype.Mat{Cols: 4, Rows: 3, Elem: gputype.F32}), "mat4x3"},
	})

	// Builtin — lifted from selena/emit/glsl/glsl.go (all identity, same as WGSL)
	run(t, []row{
		{"GLSL Builtin dot", d.Builtin("dot", []string{"a", "b"}), "dot(a, b)"},
		{"GLSL Builtin smoothstep", d.Builtin("smoothstep", []string{"e0", "e1", "x"}), "smoothstep(e0, e1, x)"},
		{"GLSL Builtin unknown", d.Builtin("myFunc", []string{"x"}), "myFunc(x)"},
	})

	// Swizzle
	run(t, []row{
		{"GLSL Swizzle xyz", d.Swizzle("v", "xyz"), "v.xyz"},
	})

	// Sample — lifted from selena/emit/glsl/glsl.go:74
	// "texture2D(" + tex + ", " + uv + ")"
	run(t, []row{
		{"GLSL Sample", d.Sample("t", "uv"), "texture2D(t, uv)"},
	})
}

// ---------------------------------------------------------------------------
// GLES
// ---------------------------------------------------------------------------

func TestGLES(t *testing.T) {
	d := dialect.GLES{}

	// TypeName — same tables as GLSL (lifted from selena/emit/gles/gles.go:76-93)
	run(t, []row{
		{"GLES TypeName float", d.TypeName(gputype.F32), "float"},
		{"GLES TypeName vec3", d.TypeName(gputype.Vec{N: 3, Elem: gputype.F32}), "vec3"},
		{"GLES TypeName vec4", d.TypeName(gputype.Vec{N: 4, Elem: gputype.F32}), "vec4"},
		{"GLES TypeName ivec2", d.TypeName(gputype.Vec{N: 2, Elem: gputype.I32}), "ivec2"},
		{"GLES TypeName uvec4", d.TypeName(gputype.Vec{N: 4, Elem: gputype.U32}), "uvec4"},
		{"GLES TypeName mat3", d.TypeName(gputype.Mat{Cols: 3, Rows: 3, Elem: gputype.F32}), "mat3"},
		{"GLES TypeName mat4", d.TypeName(gputype.Mat{Cols: 4, Rows: 4, Elem: gputype.F32}), "mat4"},
		{"GLES TypeName mat4x3", d.TypeName(gputype.Mat{Cols: 4, Rows: 3, Elem: gputype.F32}), "mat4x3"},
	})

	// Builtin
	run(t, []row{
		{"GLES Builtin dot", d.Builtin("dot", []string{"a", "b"}), "dot(a, b)"},
		{"GLES Builtin unknown", d.Builtin("myFunc", []string{"x"}), "myFunc(x)"},
	})

	// Swizzle
	run(t, []row{
		{"GLES Swizzle xyz", d.Swizzle("v", "xyz"), "v.xyz"},
	})

	// Sample — lifted from selena/emit/gles/gles.go:74
	// "texture(" + tex + ", " + uv + ")"
	run(t, []row{
		{"GLES Sample", d.Sample("t", "uv"), "texture(t, uv)"},
	})
}

// ---------------------------------------------------------------------------
// Metal
// ---------------------------------------------------------------------------

func TestMetal(t *testing.T) {
	d := dialect.Metal{}

	// TypeName — lifted from selena/emit/metal/metal.go:112-129
	// and elio/emit/metal/metal.go:316-347 for i32/u32/bool scalars and
	// Array types. Metal: float, float2, float3, float4, float3x3, float4x4.
	run(t, []row{
		{"Metal TypeName float", d.TypeName(gputype.F32), "float"},
		{"Metal TypeName int", d.TypeName(gputype.I32), "int"},
		{"Metal TypeName uint", d.TypeName(gputype.U32), "uint"},
		{"Metal TypeName bool", d.TypeName(gputype.Bool), "bool"},
		{"Metal TypeName float2", d.TypeName(gputype.Vec{N: 2, Elem: gputype.F32}), "float2"},
		{"Metal TypeName float3", d.TypeName(gputype.Vec{N: 3, Elem: gputype.F32}), "float3"},
		{"Metal TypeName float4", d.TypeName(gputype.Vec{N: 4, Elem: gputype.F32}), "float4"},
		{"Metal TypeName int3", d.TypeName(gputype.Vec{N: 3, Elem: gputype.I32}), "int3"},
		{"Metal TypeName uint4", d.TypeName(gputype.Vec{N: 4, Elem: gputype.U32}), "uint4"},
		{"Metal TypeName float3x3", d.TypeName(gputype.Mat{Cols: 3, Rows: 3, Elem: gputype.F32}), "float3x3"},
		{"Metal TypeName float4x4", d.TypeName(gputype.Mat{Cols: 4, Rows: 4, Elem: gputype.F32}), "float4x4"},
		{"Metal TypeName float4x3", d.TypeName(gputype.Mat{Cols: 4, Rows: 3, Elem: gputype.F32}), "float4x3"},
		// Array — lifted from elio/emit/metal/metal.go:339-342
		{"Metal TypeName array<f32>", d.TypeName(gputype.Array{Elem: gputype.F32, Len: 0}), "float"},
		{"Metal TypeName array<f32,6>", d.TypeName(gputype.Array{Elem: gputype.F32, Len: 6}), "float[6]"},
		{"Metal TypeName array<u32,4>", d.TypeName(gputype.Array{Elem: gputype.U32, Len: 4}), "uint[4]"},
	})

	// Builtin — Selena Metal builtinSpellings are all identity
	// (lifted from selena/emit/metal/metal.go:139-166). No renames.
	run(t, []row{
		{"Metal Builtin dot", d.Builtin("dot", []string{"a", "b"}), "dot(a, b)"},
		{"Metal Builtin mix", d.Builtin("mix", []string{"a", "b", "t"}), "mix(a, b, t)"},
		{"Metal Builtin normalize", d.Builtin("normalize", []string{"v"}), "normalize(v)"},
		{"Metal Builtin unknown", d.Builtin("myFunc", []string{"x"}), "myFunc(x)"},
	})

	// Swizzle
	run(t, []row{
		{"Metal Swizzle xyz", d.Swizzle("v", "xyz"), "v.xyz"},
	})

	// Sample — lifted from selena/emit/metal/metal.go:108-110
	// tex + ".sample(" + tex + "Sampler, " + uv + ")"
	run(t, []row{
		{"Metal Sample", d.Sample("t", "uv"), "t.sample(tSampler, uv)"},
		{"Metal Sample albedo", d.Sample("albedo", "coord"), "albedo.sample(albedoSampler, coord)"},
	})
}
