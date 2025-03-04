package core

import (
	"github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go/wamr"
)

// Context 包装WASM运行时上下文
type Context struct {
	Module      *wamr.Module
	Instance    *wamr.Instance
	Initialized bool
}

// WasmValue 表示WASM值类型
type WasmValue struct {
	Kind WasmValueKind
	Data interface{}
}

// WasmValueKind 表示WASM值的类型
type WasmValueKind int

const (
	WasmValueI32 WasmValueKind = iota
	WasmValueI64
	WasmValueF32
	WasmValueF64
)

// 导出项类型的枚举常量
const (
	ExportTypeFunc   uint8 = 0 // 函数类型
	ExportTypeGlobal uint8 = 1 // 全局变量类型
	ExportTypeMemory uint8 = 2 // 内存类型
	ExportTypeTable  uint8 = 3 // 表类型
)

// ExportType 表示导出项的类型信息
type ExportType struct {
	Kind uint8  // 导出项类型
	Name string // 导出项名称
}
