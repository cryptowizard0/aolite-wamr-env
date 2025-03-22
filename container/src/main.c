#include <lauxlib.h>
#include <lua.h>
#include <lualib.h>
#include <string.h>

#include "lsqlite3.h"

#include <stdio.h>
#include <stdlib.h>
#include <signal.h>

//#include <emscripten.h>

int boot_lua(lua_State* L);
static lua_State *wasm_lua_state = NULL;

// Pre-compiled lua loader program
static const unsigned char program[] = {__LUA_BASE__};
// Pre-compiled entry script which user wrote
static const unsigned char lua_main_program[] = {__LUA_MAIN__};


static const char* process_handle(lua_State* L, const char* json_msg, const char* json_env) {
  // 准备 Lua 代码
  const char* lua_code = 
      "print('This form lua: ')"
      "local json_msg = [=[%s]=]\n"
      "local json_env = [=[%s]=]\n"
      "print(json_msg)\n"
      "print(json_env)\n"
      "return 'hehe, ok!'";

      // "require('ao')\n"
      // "local process = require('process')\n"
      // "local json = require('json')\n"
      // "local json_msg = [=[%s]=]\n"
      // "local json_env = [=[%s]=]\n"
      // "local table_msg = json.decode(json_msg)\n"
      // "local table_env = json.decode(json_env)\n"
      // "local resp = process.handle(table_msg, table_env)\n"
      // "local json_outbox = json.encode(resp)\n"
      // "return json_outbox";

  // 分配内存用于格式化后的代码
  size_t code_size = strlen(lua_code) + strlen(json_msg) + strlen(json_env) + 1;
  char* formatted_code = (char*)malloc(code_size);
  if (!formatted_code) {
      return NULL;
  }

  // 格式化 Lua 代码
  snprintf(formatted_code, code_size, lua_code, json_msg, json_env);

  // 执行 Lua 代码
  if (luaL_dostring(L, formatted_code) != LUA_OK) {
      free(formatted_code);
      return NULL;
  }
  free(formatted_code);

  // 获取返回值
  const char* result = lua_tostring(L, -1);
  if (!result) {
      return NULL;
  }

  // 复制结果字符串（因为 Lua 可能会清理栈）
  char* final_result = strdup(result);
  lua_pop(L, 1);  // 清理栈
  fprintf(stderr, "lua handle result: %s \n", result);
  return final_result;
}

__attribute__((visibility("default")))
const char*  handle(const char* arg_0, const char* arg_1) {
  fprintf(stderr,"boot lua in handle ### arg_0: %s, arg_1: %s\n", arg_0, arg_1);
  if (wasm_lua_state == NULL) {
    fprintf(stderr, "newstate 1 \n");
    wasm_lua_state = luaL_newstate();
    fprintf(stderr,"newstate 2 \n");
    boot_lua(wasm_lua_state);
    fprintf(stderr,"newstate 3 \n");
  }
  fprintf(stderr,"get function \n");

  // call lua throw loader.lua
  //==================================
  // Push arguments
  // lua_getglobal(wasm_lua_state, "handle");
  // if (!lua_isfunction(wasm_lua_state, -1)) {
  //   printf("function handle is not defined globaly in lua runtime\n");
  //   lua_settop(wasm_lua_state, 0);
  //   return "";
  // }
  // lua_pushstring(wasm_lua_state, arg_0);
  // lua_pushstring(wasm_lua_state, arg_1);

  // // Call lua function
  // fprintf(stderr," pcall \n");
  // if (lua_pcall(wasm_lua_state, 2, 1, 0)) {
  //   printf("failed to call handle function\n");
  //   printf("error: %s\n", lua_tostring(wasm_lua_state, -1));
  //   lua_settop(wasm_lua_state, 0);
  //   return "";
  // }
  
  // // Handle return values
  // if (lua_isstring(wasm_lua_state, -1)) {
  //   const char* return_value = lua_tostring(wasm_lua_state, -1);
  //   lua_settop(wasm_lua_state, 0);
  //   return return_value;
  // }
  // return "";
  //==================================

  // call lua throw do_string
  return process_handle(wasm_lua_state, arg_0, arg_1);
  
}

// This function is for debug to see an C <-> Lua stack values
// void dumpStack(lua_State *L) {
//   int i;
//   int stackSize = lua_gettop(L);
//   for (i = stackSize; i >= 1; i--) {
//     int stackType = lua_type(L, i);
//     printf("Stack[%2d-%10s]:", i, lua_typename(L, stackType));
// 
//     switch (stackType) {
//       case LUA_TNUMBER:
//         printf("%f", lua_tonumber(L, i));
//         break;
//       case LUA_TBOOLEAN:
//         if (lua_toboolean(L, i)) {
//           printf("true");
//         } else {
//           printf("false");
//         }
//         break;
//       case LUA_TSTRING:
//         printf("%s", lua_tostring(L, i));
//         break;
//       case LUA_TNIL:
//         printf("nil");
//         break;
//       default:
//         printf("%s", lua_typename(L, stackType));
//         break;
//     }
//     printf("\n");
//   }
//   printf("\n");
// }

/* Copied from lua.c */

static lua_State *globalL = NULL;

static void lstop (lua_State *L, lua_Debug *ar) {
  (void)ar;  /* unused arg. */
  lua_sethook(L, NULL, 0, 0);  /* reset hook */
  luaL_error(L, "interrupted!");
}

static void laction (int i) {
  signal(i, SIG_DFL); /* if another SIGINT happens, terminate process */
  lua_sethook(globalL, lstop, LUA_MASKCALL | LUA_MASKRET | LUA_MASKCOUNT, 1);
}

static int msghandler (lua_State *L) {
  const char *msg = lua_tostring(L, 1);
  if (msg == NULL) {  /* is error object not a string? */
    if (luaL_callmeta(L, 1, "__tostring") &&  /* does it have a metamethod */
        lua_type(L, -1) == LUA_TSTRING)  /* that produces a string? */
      return 1;  /* that is the message */
    else
      msg = lua_pushfstring(L, "(error object is a %s value)",
                               luaL_typename(L, 1));
  }
  /* Call debug.traceback() instead of luaL_traceback() for Lua 5.1 compatibility. */
  lua_getglobal(L, "debug");
  lua_getfield(L, -1, "traceback");
  /* debug */
  lua_remove(L, -2);
  lua_pushstring(L, msg);
  /* original msg */
  lua_remove(L, -3);
  lua_pushinteger(L, 2);  /* skip this function and traceback */
  lua_call(L, 2, 1); /* call debug.traceback */
  return 1;  /* return the traceback */
}

static int docall (lua_State *L, int narg, int nres) {
  int status;
  int base = lua_gettop(L) - narg;  /* function index */
  lua_pushcfunction(L, msghandler);  /* push message handler */
  lua_insert(L, base);  /* put it under function and args */
  globalL = L;  /* to be available to 'laction' */
  signal(SIGINT, laction);  /* set C-signal handler */
  status = lua_pcall(L, narg, nres, base);
  signal(SIGINT, SIG_DFL); /* reset C-signal handler */
  lua_remove(L, base);  /* remove message handler from the stack */
  return status;
}

// Boot function
int main(void) {
  volatile int debug = 1;
  if (wasm_lua_state != NULL) {
    return 0;
  }
  wasm_lua_state = luaL_newstate();
  if (boot_lua(wasm_lua_state)) {
    printf("failed to boot lua runtime\\n");
    lua_close(wasm_lua_state);
    return 1;
  }
  if (debug) printf("Boot Lua Webassembly!\n");
  // test lua print
  lua_getglobal(wasm_lua_state, "print");
  lua_pushstring(wasm_lua_state, "Hello, Lua from main!");
  if (lua_pcall(wasm_lua_state, 1, 0, 0) != 0) {
    printf("error running function `print`: %s\n", lua_tostring(wasm_lua_state, -1));
  }
  // const char* env = "{__ENV__}";
  // const char* msg = "{__MSG__}";
  // if (debug) printf("call handle(%s, %s) \n", msg, env);
  
  // const char* result = handle(msg, env);
  // printf("result: %s\n", result);
  
  return 0;
}

// boot lua runtime from compiled lua source
int boot_lua(lua_State* L) {
  luaL_openlibs(L);
  printf("checkstack!\n");
  if (!lua_checkstack(L, 8192)) {
    fprintf(stderr, "无法分配足够的栈空间\n");
    return 1;
  }
  
  // Preload lsqlite3
  luaL_getsubtable(L, LUA_REGISTRYINDEX, LUA_PRELOAD_TABLE);
  lua_pushcfunction(L, luaopen_lsqlite3);
  lua_setfield(L, -2, LUA_SQLLIBNAME);
  lua_pop(L, 1);  // remove PRELOAD table

  if (luaL_loadbuffer(L, (const char*)program, sizeof(program), "main")) {
    fprintf(stderr, "error on luaL_loadbuffer()\n");
    return 1;
  }
  lua_newtable(L);
  lua_pushlstring(L, (const char*)lua_main_program, sizeof(lua_main_program));
  lua_setfield(L, -2, "__lua_webassembly__");

  // This place will be injected by emcc-lua
  __INJECT_LUA_FILES__

  if (docall(L, 1, LUA_MULTRET)) {
    const char *errmsg = lua_tostring(L, 1);
    if (errmsg) {
      fprintf(stderr, "%s\n", errmsg);
    }
    lua_close(L);
    return 1;
  }
  return 0;
}