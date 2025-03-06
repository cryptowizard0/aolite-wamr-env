package core

const (
	WasmValueI32       = 1
	WasmValueI64       = 2
	WasmValueF32       = 3
	WasmValueF64       = 4
	WasmValueAnyRef    = 5
	WasmValueFuncRef   = 6
	WasmValueExternRef = 7
	WasmValueString    = 8 // 添加字符串类型
)

// 导出项类型的枚举常量
const (
	ExportTypeFunc   uint8 = 0 // 函数类型
	ExportTypeGlobal uint8 = 1 // 全局变量类型
	ExportTypeMemory uint8 = 2 // 内存类型
	ExportTypeTable  uint8 = 3 // 表类型
)

// WasmValue 表示WASM值类型
type WasmValue struct {
	Kind WasmValueKind
	Data interface{}
}

// WasmValueKind 表示WASM值的类型
type WasmValueKind int

// ExportType 表示导出项的类型信息
type ExportType struct {
	Kind uint8  // 导出项类型
	Name string // 导出项名称
}
