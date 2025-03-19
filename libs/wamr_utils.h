#ifndef WAMR_UTILS_H
#define WAMR_UTILS_H

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <wasm_export.h>
#include <wasm_c_api.h>

#ifdef __cplusplus
extern "C" {
#endif

// 导出函数声明
extern void print_imports(wasm_module_t module);
extern void print_imports2(wasm_module_t module, wasm_module_inst_t inst);
extern bool register_wasi_clock();
extern bool register_tests();
extern bool register_env();

#ifdef __cplusplus
}
#endif

#endif // WAMR_UTILS_H