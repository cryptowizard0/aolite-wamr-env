#include <setjmp.h>
#include "wamr_imports.h"


// global
static jmp_buf g_jmpbuf;
static int32_t g_tempRet0 = 0;


void print_imports2(wasm_module_t module, wasm_module_inst_t inst) {
    uint32_t import_count = wasm_runtime_get_import_count(module);
    printf("Total imports: %d\n", import_count);

    for (uint32_t i = 0; i < import_count; i++) {
        wasm_import_t import_type;
        wasm_runtime_get_import_type(module, i, &import_type);
        void *func_ptr = wasm_runtime_lookup_function(inst, import_type.name);

        printf("Import #%d:\n", i);
        printf("  Module name: %s\n", import_type.module_name);
        printf("  Name: %s\n", import_type.name);
        printf("  Kind: %d\n", import_type.kind);
        printf("  Linked: %s\n\n", func_ptr ? "Yes" : "No");
    }
}

// WASI clock_time_get 实现
static int32_t wasi_clock_time_get(wasm_exec_env_t exec_env, int32_t clock_id,
                                  int64_t precision, int32_t time_ptr) {
    wasm_module_inst_t module_inst = get_module_inst(exec_env);
    if (!module_inst)
        return -1;

    if (!wasm_runtime_validate_app_addr(module_inst, time_ptr, sizeof(int64_t)))
        return -1;

    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts) != 0)
        return -1;

    int64_t time_ns = (int64_t)ts.tv_sec * 1000000000LL + ts.tv_nsec;
    void *native_ptr = wasm_runtime_addr_app_to_native(module_inst, time_ptr);
    *(int64_t*)native_ptr = time_ns;

    return 0;
}


static uint32_t env_signal(wasm_exec_env_t exec_env, uint32_t signum, uint32_t handler) {
    printf("Signal called with signum: %d, handler: %d\n", signum, handler);
    return 0;  // Return 0 for success
}

uint32_t env_saveSetjmp(wasm_exec_env_t exec_env, uint32_t env_addr, uint32_t label, uint32_t table, uint32_t size) {
    wasm_module_inst_t module_inst = get_module_inst(exec_env);
    if (!module_inst)
        return -1;

    if (!wasm_runtime_validate_app_addr(module_inst, env_addr, sizeof(int64_t)))
        return -1;

    return setjmp(g_jmpbuf);
}

uint32_t env_wasm_setjmp_test(wasm_exec_env_t exec_env) {
    return setjmp(g_jmpbuf);
}

void env_emscripten_longjmp(wasm_exec_env_t exec_env, uint32_t env, uint32_t val) {
    longjmp(g_jmpbuf, 1);
}

void env_setTempRet0(wasm_exec_env_t exec_env, int32_t value) {
    g_tempRet0 = value;
}

uint32_t env_getTempRet0(wasm_exec_env_t exec_env) {
    return g_tempRet0;
}

void print_imports(wasm_module_t module) {
    uint32_t import_count = wasm_runtime_get_import_count(module);
    printf("Total imports: %d\n", import_count);

    for (uint32_t i = 0; i < import_count; i++) {
        wasm_import_t import_type;
        wasm_runtime_get_import_type(module, i, &import_type);

        printf("Import #%d:\n", i);
        printf("  Module name: %s\n", import_type.module_name);
        printf("  Name: %s\n", import_type.name);
        printf("  Kind: %d\n\n", import_type.kind);
    }
}


// 其他函数实现...
static uint32_t mul(wasm_exec_env_t exec_env, uint32_t a, uint32_t b) { 
    return a * b; 
}



static NativeSymbol native_symbols_test[] = {
    {"mul", mul, "(ii)i", NULL}
};

static NativeSymbol native_symbols_env[] = {
    {"saveSetjmp", env_saveSetjmp, "(iiii)i", NULL},
    {"__wasm_longjmp", env_emscripten_longjmp, "(ii)", NULL}, 
    {"getTempRet0", env_getTempRet0, "()i", NULL},
    {"signal", env_signal, "(ii)i", NULL},
};

// 注册函数实现
bool register_tests() {
    printf("Import tests \n");
    return wasm_runtime_register_natives("my_lib", native_symbols_test, 1);
}

// ... existing code ...

bool register_env() {
    printf("\n=== Registering env functions ===\n");
    size_t func_count = sizeof(native_symbols_env) / sizeof(native_symbols_env[0]);
    printf("Registering %zu functions for 'env' module:\n", func_count);
    
    for (size_t i = 0; i < func_count; i++) {
        printf("Function[%zu]:\n", i);
        printf("  Name: %s\n", native_symbols_env[i].symbol);
        printf("  Signature: %s\n", native_symbols_env[i].signature);
    }
    
    bool result = wasm_runtime_register_natives("env", native_symbols_env, func_count);
    printf("Registration %s\n", result ? "succeeded" : "failed");
    
    if (!result) {
        printf("!!! Registration failed for 'env' module !!!\n");
        // 尝试单独注册每个函数以定位问题
        for (size_t i = 0; i < func_count; i++) {
            bool single_result = wasm_runtime_register_natives(
                "env", &native_symbols_env[i], 1);
            printf("Individual registration of '%s': %s\n", 
                   native_symbols_env[i].symbol,
                   single_result ? "succeeded" : "failed");
        }
    }
    
    printf("===============================\n\n");
    return result;
}