package dialect

import (
	"fmt"
	"strings"

	"m31labs.dev/prism/gputype"
)

// GLSL implements [Dialect] for GLSL ES 1.00 (WebGL / desktop GL).
// TypeName spellings lifted from selena/emit/glsl/glsl.go:76-93 and
// elio/emit/glsl/glsl.go:442-480 (i32/u32 vectors, non-square mats).
// Builtin spellings lifted from selena/emit/glsl/glsl.go:103-130.
// Sample lifted from selena/emit/glsl/glsl.go:74.
type GLSL struct{}

func (GLSL) TypeName(t gputype.Type) string { return glslTypeName(t) }

// Builtin renders a builtin call. Unknown names fall back to the canonical name.
// Lifted from selena/emit/internal/spell/spell.go:10-16.
func (GLSL) Builtin(name string, args []string) string {
	if spelled, ok := glslBuiltins[name]; ok {
		name = spelled
	}
	return name + "(" + strings.Join(args, ", ") + ")"
}

func (GLSL) Swizzle(expr, components string) string { return expr + "." + components }

// Sample renders: texture2D(tex, uv)
// Lifted from selena/emit/glsl/glsl.go:74.
func (GLSL) Sample(tex, uv string) string { return "texture2D(" + tex + ", " + uv + ")" }

// glslTypeName is shared by both GLSL and GLES (same type table).
// Scalars: lifted from elio/emit/glsl/glsl.go:444-451.
// Vectors: lifted from elio/emit/glsl/glsl.go:452-462.
// Mats: lifted from elio/emit/glsl/glsl.go:463-468.
func glslTypeName(t gputype.Type) string {
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
			return x.Name
		}
	case gputype.Vec:
		switch x.Elem.Name {
		case "u32":
			return fmt.Sprintf("uvec%d", x.N)
		case "i32":
			return fmt.Sprintf("ivec%d", x.N)
		default:
			return fmt.Sprintf("vec%d", x.N)
		}
	case gputype.Mat:
		if x.Cols == x.Rows {
			return fmt.Sprintf("mat%d", x.Cols)
		}
		return fmt.Sprintf("mat%dx%d", x.Cols, x.Rows)
	case gputype.Array:
		// GLSL puts array size after the name at the declaration site;
		// TypeName returns the element type (callers append [N] themselves).
		if x.Len == 0 {
			return glslTypeName(x.Elem)
		}
		return fmt.Sprintf("%s[%d]", glslTypeName(x.Elem), x.Len)
	}
	return "/* unknown type */"
}

// glslBuiltins maps canonical builtin names to their GLSL spelling.
// Lifted from selena/emit/glsl/glsl.go:103-130.
var glslBuiltins = map[string]string{
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
