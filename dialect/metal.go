package dialect

import (
	"fmt"
	"strings"

	"m31labs.dev/prism/gputype"
)

// Metal implements [Dialect] for Metal Shading Language (iOS / macOS).
// TypeName spellings lifted from selena/emit/metal/metal.go:112-129 and
// elio/emit/metal/metal.go:316-347 (i32/u32/bool scalars, Array).
// Builtin spellings lifted from selena/emit/metal/metal.go:139-166 (all identity).
// Sample lifted from selena/emit/metal/metal.go:108-110.
type Metal struct{}

func (Metal) TypeName(t gputype.Type) string { return metalTypeName(t) }

// Builtin renders a builtin call. Unknown names fall back to the canonical name.
// Lifted from selena/emit/internal/spell/spell.go:10-16.
// Note: Selena's Metal builtins map has all-identity spellings (no renames).
func (Metal) Builtin(name string, args []string) string {
	if spelled, ok := metalBuiltins[name]; ok {
		name = spelled
	}
	return name + "(" + strings.Join(args, ", ") + ")"
}

func (Metal) Swizzle(expr, components string) string { return expr + "." + components }

// Sample renders: tex.sample(texSampler, uv)
// Lifted from selena/emit/metal/metal.go:108-110.
func (Metal) Sample(tex, uv string) string {
	return tex + ".sample(" + tex + "Sampler, " + uv + ")"
}

// metalTypeName spells a gputype.Type in Metal Shading Language.
// Scalar mapping lifted from elio/emit/metal/metal.go:318-325:
//
//	f32 → float, u32 → uint, i32 → int, bool → bool, else pass-through.
//
// Vec lifted from elio/emit/metal/metal.go:329-330: "<elemType><N>", e.g. float3.
// Mat lifted from elio/emit/metal/metal.go:331-332: "<elemType><Cols>x<Rows>", e.g. float4x4.
// Array lifted from elio/emit/metal/metal.go:337-342: runtime (Len==0) → elem type;
// fixed (Len>0) → "<elemType>[N]".
func metalTypeName(t gputype.Type) string {
	switch x := t.(type) {
	case gputype.Scalar:
		switch x.Name {
		case "f32":
			return "float"
		case "u32":
			return "uint"
		case "i32":
			return "int"
		default:
			return x.Name // bool, etc.
		}
	case gputype.Vec:
		return fmt.Sprintf("%s%d", metalTypeName(x.Elem), x.N)
	case gputype.Mat:
		return fmt.Sprintf("%s%dx%d", metalTypeName(x.Elem), x.Cols, x.Rows)
	case gputype.Array:
		if x.Len == 0 {
			return metalTypeName(x.Elem)
		}
		return fmt.Sprintf("%s[%d]", metalTypeName(x.Elem), x.Len)
	}
	return "/* unknown type */"
}

// metalBuiltins maps canonical builtin names to their Metal spelling.
// Lifted from selena/emit/metal/metal.go:139-166: all identity (no renames).
var metalBuiltins = map[string]string{
	"abs":        "abs",
	"clamp":      "clamp",
	"distance":   "distance",
	"dot":        "dot",
	"length":     "length",
	"max":        "max",
	"min":        "min",
	"mix":        "mix",
	"normalize":  "normalize",
	"pow":        "pow",
	"sin":        "sin",
	"cos":        "cos",
	"tan":        "tan",
	"sqrt":       "sqrt",
	"floor":      "floor",
	"ceil":       "ceil",
	"fract":      "fract",
	"sign":       "sign",
	"exp":        "exp",
	"log":        "log",
	"exp2":       "exp2",
	"log2":       "log2",
	"cross":      "cross",
	"reflect":    "reflect",
	"step":       "step",
	"smoothstep": "smoothstep",
}
