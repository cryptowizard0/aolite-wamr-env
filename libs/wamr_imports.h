#ifndef WAMR_IMPORTS_H
#define WAMR_IMPORTS_H

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

// stdlib
extern int32_t env_wasm_setjmp(wasm_exec_env_t exec_env, int32_t buf, int32_t label, int32_t table);
extern int32_t env_wasm_setjmp_test(wasm_exec_env_t exec_env);
extern void env_emscripten_longjmp(wasm_exec_env_t exec_env);
extern void env_setTempRet0(wasm_exec_env_t exec_env, int32_t value);
extern int32_t env_getTempRet0(wasm_exec_env_t exec_env);

#ifdef __cplusplus
}
#endif

#endif // WAMR_UTILS_H