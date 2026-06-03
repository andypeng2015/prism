package dialect

import (
	"fmt"
	"strings"

	"m31labs.dev/prism/gputype"
)

// WGSL implements [Dialect] for the WebGPU Shading Language.
// TypeName spellings lifted from selena/emit/wgsl/wgsl.go:110-127 and
// elio/emit/wgsl/wgsl.go:199-217 (Array).
// Builtin spellings lifted from selena/emit/wgsl/wgsl.go:137-164.
// Sample lifted from selena/emit/wgsl/wgsl.go:106-108.
type WGSL struct{}

func (WGSL) TypeName(t gputype.Type) string {
	switch x := t.(type) {
	case gputype.Scalar:
		return x.Name // f32, i32, u32, bool
	case gputype.Vec:
		return fmt.Sprintf("vec%d<%s>", x.N, x.Elem.Name)
	case gputype.Mat:
		return fmt.Sprintf("mat%dx%d<%s>", x.Cols, x.Rows, x.Elem.Name)
	case gputype.Array:
		if x.Len == 0 {
			return fmt.Sprintf("array<%s>", WGSL{}.TypeName(x.Elem))
		}
		return fmt.Sprintf("array<%s, %d>", WGSL{}.TypeName(x.Elem), x.Len)
	}
	return "/* unknown type */"
}

// Builtin renders a builtin call. Unknown names fall back to the canonical name
// (lifted from selena/emit/internal/spell/spell.go:10-16).
func (WGSL) Builtin(name string, args []string) string {
	if spelled, ok := wgslBuiltins[name]; ok {
		name = spelled
	}
	return name + "(" + strings.Join(args, ", ") + ")"
}

func (WGSL) Swizzle(expr, components string) string { return expr + "." + components }

// Sample renders: textureSample(tex, texSampler, uv)
// Lifted from selena/emit/wgsl/wgsl.go:106-108.
func (WGSL) Sample(tex, uv string) string {
	return "textureSample(" + tex + ", " + tex + "Sampler, " + uv + ")"
}

// wgslBuiltins maps canonical builtin names to their WGSL spelling.
// Lifted from selena/emit/wgsl/wgsl.go:137-164.
var wgslBuiltins = map[string]string{
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
