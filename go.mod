module aolite-wamr-evn

go 1.23.4

require github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go v0.0.0-00010101000000-000000000000

replace github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go => ./_build/wamr/language-bindings/go
