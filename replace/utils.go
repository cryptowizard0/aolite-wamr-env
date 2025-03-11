package wamr

/*
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <wasm_export.h>
#include <wasm_c_api.h>

// 定义一个函数来打印导入表
void print_imports(wasm_module_t module) {
    uint32_t import_count = wasm_runtime_get_import_count(module);
    printf("Total imports: %d\n", import_count);

    // wasm_importtype_vec_t import_info;
    // const wasm_module_t* module_ptr = &module;
    // wasm_module_imports(module_ptr, &import_info);
    for (uint32_t i = 0; i < import_count; i++) {
        wasm_import_t import_type;

        wasm_runtime_get_import_type(module, i, &import_type);

		// void *func_ptr = wasm_runtime_lookup_function(inst, import_type.name);

        printf("Import #%d:\n", i);
        printf("  Module name: %s\n", import_type.module_name);
        printf("  Name: %s\n", import_type.name);  // field_name 改为 name
        printf("  Kind: %d\n\n", import_type.kind);
		// printf("  Linked: %s\n\n", func_ptr ? "Yes" : "No");
    }
}

void print_imports2(wasm_module_t module, wasm_module_inst_t inst) {
    uint32_t import_count = wasm_runtime_get_import_count(module);
    printf("Total imports: %d\n", import_count);

    // wasm_importtype_vec_t import_info;
    // const wasm_module_t* module_ptr = &module;
    // wasm_module_imports(module_ptr, &import_info);
    for (uint32_t i = 0; i < import_count; i++) {
        wasm_import_t import_type;

        wasm_runtime_get_import_type(module, i, &import_type);

		void *func_ptr = wasm_runtime_lookup_function(inst, import_type.name);

        printf("Import #%d:\n", i);
        printf("  Module name: %s\n", import_type.module_name);
        printf("  Name: %s\n", import_type.name);  // field_name 改为 name
        printf("  Kind: %d\n", import_type.kind);
		printf("  Linked: %s\n\n", func_ptr ? "Yes" : "No");
    }
}


// 实现 clock_time_get（简化为示例）
// int32_t clock_time_get(int32_t clock_id, int64_t precision) {
//     // 这里应调用系统时间函数，例如 gettimeofday
//     // *time = 1234567; // 模拟时间戳
//     return 0; // 成功返回 0
// }

// WASI clock_time_get 实现
static int32_t wasi_clock_time_get(wasm_exec_env_t exec_env, int32_t clock_id,
                                  int64_t precision, int32_t time_ptr) {
    // 获取 WASM 实例
    wasm_module_inst_t module_inst = get_module_inst(exec_env);
    if (!module_inst)
        return -1;

    // 验证内存访问
    if (!wasm_runtime_validate_app_addr(module_inst, time_ptr, sizeof(int64_t)))
        return -1;

    // 获取当前时间
    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts) != 0)
        return -1;

    // 转换为纳秒
    int64_t time_ns = (int64_t)ts.tv_sec * 1000000000LL + ts.tv_nsec;

    // 写入 WASM 内存
    void *native_ptr = wasm_runtime_addr_app_to_native(module_inst, time_ptr);
    *(int64_t*)native_ptr = time_ns;

    return 0;
}

uint32_t mul(wasm_exec_env_t exec_env, uint32_t a, uint32_t b) { return a * b; }

static NativeSymbol native_symbols[] = {
    {
        "clock_time_get",           // 函数名
        wasi_clock_time_get,        // 实现函数
        "(iI*)i",                   // 函数签名
        NULL                        // 附加数据
    }
};

static NativeSymbol native_symbols2[] = {
	{"mul", mul, "(ii)i", NULL}
    // 其他函数...
};


// 注册函数到 WAMR
bool register_wasi_clock() {
    return wasm_runtime_register_natives("wasi_snapshot_preview1",
                                       native_symbols,
                                    1);
}

bool register_tests() {
    return wasm_runtime_register_natives("my_lib", native_symbols2, 1);
}
*/
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
	success := C.register_wasi_clock()
	if success {
		fmt.Println("WASI functions registered successfully")
	} else {
		fmt.Println("Failed to register WASI functions")
		return fmt.Errorf("Failed to register WASI functions")
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
func (self *Instance) WriteString(s string) (int64, error) {
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

	return int64(uintptr(ptr)), nil
}

// ReadString 从 WASM 内存中读取字符串
func (self *Instance) ReadString(ptr int64) string {
	if ptr == 0 {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(uintptr(ptr))))
}

// FreeString 释放字符串占用的内存
func (self *Instance) FreeString(ptr int64) {
	if ptr != 0 {
		C.wasm_runtime_free(unsafe.Pointer(uintptr(ptr)))
	}
}
