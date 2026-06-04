// Package prism is the shared GPU-codegen toolkit for M31 Labs' GPU-DSL
// compilers (Elio, Selena, Eos) and gosx.
//
// Prism provides three things:
//
//   - A canonical GPU type vocabulary — scalar, vector, matrix, sampler, and
//     texture types that every compiler can reference without re-declaring them.
//     Consumers opt in; nothing is forced on the host AST.
//
//   - A Dialect layer that spells types, builtins, and swizzles correctly for
//     each shader backend: WGSL, GLSL (desktop/ES), and Metal (MSL). Compilers
//     delegate spelling decisions to Prism rather than hard-coding per-backend
//     string tables.
//
//   - A validation harness that invokes offline validators — naga (for WGSL) and
//     glslangValidator (for GLSL/GLES) — and surfaces structured diagnostics back
//     to the caller.
//
//   - Backend source descriptors and entry-shape checks for generated GPU kernel
//     variants. This lets compilers such as Eos validate emitted CUDA, Metal,
//     Vulkan, DirectML, and WebGPU source without Prism knowing their artifact IR.
//
// Prism depends only on gotreesitter. It is consumed by Elio, Selena, Eos, and
// gosx, and it deliberately never imports any of those packages.
package prism
