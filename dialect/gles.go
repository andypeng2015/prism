package dialect

import (
	"strings"

	"m31labs.dev/prism/gputype"
)

// GLES implements [Dialect] for GLSL ES 3.00 (Android / GLES 3.1).
// TypeName spellings lifted from selena/emit/gles/gles.go:76-93 — identical
// to GLSL; the difference is only the Sample call form and shader IO keywords.
// Builtin spellings lifted from selena/emit/gles/gles.go:103-130.
// Sample lifted from selena/emit/gles/gles.go:74.
type GLES struct{}

func (GLES) TypeName(t gputype.Type) string { return glslTypeName(t) }

// Builtin renders a builtin call. Unknown names fall back to the canonical name.
// Lifted from selena/emit/internal/spell/spell.go:10-16.
func (GLES) Builtin(name string, args []string) string {
	if spelled, ok := glesBuiltins[name]; ok {
		name = spelled
	}
	return name + "(" + strings.Join(args, ", ") + ")"
}

func (GLES) Swizzle(expr, components string) string { return expr + "." + components }

// Sample renders: texture(tex, uv)
// Lifted from selena/emit/gles/gles.go:74.
func (GLES) Sample(tex, uv string) string { return "texture(" + tex + ", " + uv + ")" }

// glesBuiltins maps canonical builtin names to their GLSL ES 3.00 spelling.
// Lifted from selena/emit/gles/gles.go:103-130.
var glesBuiltins = map[string]string{
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
