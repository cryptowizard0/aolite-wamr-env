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
	// malloc memory for string, (size = len(s) + 1 for '\n' terminator in C)
	size := len(s) + 1
	offset, native_addr := self.ModuleMalloc(uint64(size))
	if native_addr == nil {
		return 0, fmt.Errorf("failed to allocate WASM memory")
	}

	// copy string from native to WASM memory
	strBytes := []byte(s + string(rune(0)))
	C.memcpy(unsafe.Pointer(native_addr), unsafe.Pointer(&strBytes[0]), C.size_t(len(strBytes)))
	return int32(offset), nil
}

// Read string from WASM memory
func (self *Instance) ReadString(ptr int32) (string, error) {
	if ptr == 0 {
		return "", nil
	}
	if !self.ValidateStrAddr(uint64(ptr)) {
		errMsg := fmt.Sprintf("read string failed: validate addrress failed:, %v", ptr)
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	// Convert WASM offset to native pointer
	native_ptr := self.AddrAppToNative(uint64(ptr))
	if native_ptr == nil {
		errMsg := fmt.Sprintf("read string failed: AddrAppToNative failed:, %v", ptr)
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	return C.GoString((*C.char)(unsafe.Pointer(native_ptr))), nil
}

// FreeString frees memory allocated in WASM for a string
func (self *Instance) FreeString(ptr int32) error {
	if ptr == 0 {
		return nil
	}
	if !self.ValidateStrAddr(uint64(ptr)) {
		errMsg := fmt.Sprintf("FreeString failed: validate addrress failed:, %v", ptr)
		return fmt.Errorf("%s", errMsg)
	}

	// Free memory
	self.ModuleFree(uint64(ptr))
	return nil
}
