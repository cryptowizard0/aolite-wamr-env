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




int32_t env_wasm_setjmp(wasm_exec_env_t exec_env, int32_t buf, int32_t label, int32_t table) {
    wasm_module_inst_t module_inst = get_module_inst(exec_env);
    if (!module_inst)
        return -1;

    if (!wasm_runtime_validate_app_addr(module_inst, buf, sizeof(int64_t)))
        return -1;

    return setjmp(g_jmpbuf);
}

int32_t env_wasm_setjmp_test(wasm_exec_env_t exec_env) {
    return setjmp(g_jmpbuf);
}

void env_emscripten_longjmp(wasm_exec_env_t exec_env) {
    longjmp(g_jmpbuf, 1);
}

void env_setTempRet0(wasm_exec_env_t exec_env, int32_t value) {
    g_tempRet0 = value;
}

int32_t env_getTempRet0(wasm_exec_env_t exec_env) {
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

static void env_invoke_vii(wasm_exec_env_t exec_env, int32_t func_index, int32_t arg1, int32_t arg2) {
    wasm_module_inst_t module_inst = wasm_runtime_get_module_inst(exec_env);
    char func_name[32];
    snprintf(func_name, sizeof(func_name), "func_%d", func_index);
    wasm_function_inst_t func = wasm_runtime_lookup_function(module_inst, func_name);

    if (func == NULL) {
        printf("Function not found\n");
        return;
    }

    uint32_t argv[2] = {arg1, arg2};
    if (!wasm_runtime_call_wasm(exec_env, func, 2, argv)) {
        printf("Failed to call function\n");
        return;
    }
}



static NativeSymbol native_symbols2[] = {
    {"mul", mul, "(ii)i", NULL}
};

static NativeSymbol native_symbols_env[] = {
    {"invoke_vii", (void*)env_invoke_vii, "(iii)", NULL},
    {"__wasm_setjmp", (void*)env_wasm_setjmp, "()i", NULL},
    {"__wasm_setjmp_test", (void*)env_wasm_setjmp_test, "()i", NULL},
    {"emscripten_longjmp", (void*)env_emscripten_longjmp, "()", NULL},
    {"setTempRet0", (void*)env_setTempRet0, "(i)", NULL},
    {"getTempRet0", (void*)env_getTempRet0, "()i", NULL}
};

// 注册函数实现
bool register_tests() {
    return wasm_runtime_register_natives("my_lib", native_symbols2, 1);
}

bool register_env() {
    return wasm_runtime_register_natives("env", native_symbols_env, sizeof(native_symbols_env) / sizeof(native_symbols_env[0]));
}