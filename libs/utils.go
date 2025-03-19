package wamr

// #cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/packaged/include
// #cgo LDFLAGS: ${SRCDIR}/wamr_utils.o
// #include "wamr_utils.h"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// ExportType 表示导出项的类型信息
type ExportType struct {
	Kind uint8
	Name string
	// 其他字段...
}

func RegisterImportFunctions() error {
	// success := C.register_wasi_clock()
	// if success {
	// 	fmt.Println("WASI functions registered successfully")
	// } else {
	// 	fmt.Println("Failed to register WASI functions")
	// 	return fmt.Errorf("Failed to register WASI functions")
	// }

	success := C.register_env()
	if success {
		fmt.Println("env functions registered successfully")
	} else {
		fmt.Println("Failed to register env functions")
		return fmt.Errorf("Failed to register env functions")
	}

	// success = C.register_tests()
	// if success {
	// 	fmt.Println("tests functions registered successfully")
	// } else {
	// 	fmt.Println("Failed to register tests functions")
	// 	return fmt.Errorf("Failed to register tests functions")
	// }

	return nil
}

func (self *Module) PrintImports() {
	C.print_imports(self.module)
}

// GetExportCount 获取模块导出项的数量
func (self *Module) GetExportCount() int32 {
	return int32(C.wasm_runtime_get_export_count(self.module))
}

// GetExportType 获取指定索引的导出项信息
func (self *Module) GetExportType(index int32) *ExportType {
	var cExport C.wasm_export_t
	C.wasm_runtime_get_export_type(self.module, C.int32_t(index), &cExport)

	return &ExportType{
		Kind: uint8(cExport.kind),
		Name: C.GoString(cExport.name),
		// 转换其他字段...
	}
}

func (self *Instance) PrintImports() {
	C.print_imports2(self._module.module, self._instance)
}

// WriteString 将字符串写入 WASM 内存并返回指针
func (self *Instance) WriteString(s string) (int32, error) {
	// 分配内存（字符串长度 + 1 用于 null 终止符）
	size := len(s) + 1
	ptr := C.wasm_runtime_malloc(C.uint(size))
	if ptr == nil {
		return 0, errors.New("failed to allocate memory")
	}

	// 将字符串复制到 WASM 内存
	data := C.GoBytes(unsafe.Pointer(ptr), C.int(size))
	copy(data, s)
	data[len(s)] = 0 // null 终止符

	return int32(uintptr(ptr)), nil
}

// ReadString 从 WASM 内存中读取字符串
func (self *Instance) ReadString(ptr int32) string {
	if ptr == 0 {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(uintptr(ptr))))
}

// FreeString 释放字符串占用的内存
func (self *Instance) FreeString(ptr int32) {
	if ptr != 0 {
		C.wasm_runtime_free(unsafe.Pointer(uintptr(ptr)))
	}
}
