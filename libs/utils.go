package wamr

// #cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/packaged/include
// #cgo LDFLAGS: ${SRCDIR}/libwamr_imports.a
// #include "wamr_imports.h"
import "C"
import (
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
	fmt.Println("call register_env")
	success := C.register_env()
	if success {
		fmt.Println("env functions registered successfully")
	} else {
		fmt.Println("Failed to register env functions")
		return fmt.Errorf("Failed to register env functions")
	}

	success = C.register_tests()
	if success {
		fmt.Println("tests functions registered successfully")
	} else {
		fmt.Println("Failed to register tests functions")
		return fmt.Errorf("Failed to register tests functions")
	}

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
	// var offset uint32
	// ptr := C.wasm_runtime_module_malloc(self._instance, C.uint64_t(size), (*unsafe.Pointer)(unsafe.Pointer(&offset)))
	// if ptr == nil {
	// 	return 0, errors.New("failed to allocate memory")
	// }
	fmt.Println("--> 1")
	offset, native_addr := self.ModuleMalloc(uint64(size))
	if native_addr == nil {
		fmt.Println("--> 2.0")
		return 0, fmt.Errorf("failed to allocate memory")
	}

	fmt.Println("--> 2", offset, native_addr)
	v := self.ValidateStrAddr(offset)
	if !v {
		fmt.Println("--> 2.1")
		return 0, fmt.Errorf("failed to validate app addr")
	}
	// 获取本地内存指针
	// native_ptr := C.wasm_runtime_addr_app_to_native(self._instance, C.uint32_t(offset))
	// if native_ptr == nil {
	// 	C.wasm_runtime_module_free(self._instance, C.uint32_t(offset))
	// 	return 0, errors.New("failed to get native pointer")
	// }
	native_ptr := self.AddrAppToNative(uint64(offset))
	fmt.Println("--> 3")
	// 将字符串复制到 WASM 内存
	strBytes := []byte(s + string(rune(0)))
	C.memcpy(unsafe.Pointer(native_ptr), unsafe.Pointer(&strBytes[0]), C.size_t(len(strBytes)))
	fmt.Println("--> 4")
	return int32(offset), nil
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
