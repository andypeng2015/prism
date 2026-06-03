# prism

**The shared GPU-codegen substrate for M31 Labs' GPU-DSL compilers.** Prism provides a canonical type vocabulary, a Dialect layer that spells types and builtins correctly for each shader backend (WGSL, GLSL, GLSL-ES, Metal), and a validation harness that runs naga and glslangValidator against emitted shaders. It is consumed by Elio, Selena, Eos, and gosx, and it never depends on any of them.
