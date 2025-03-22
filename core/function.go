package core

import (
	"errors"
	"fmt"
	"unsafe"
)

// CallFunction 调用WASM函数
func (c *Context) CallFunction(funcName string, args []WasmValue) ([]WasmValue, error) {
	if !c.Initialized {
		return nil, errors.New("runtime not initialized")
	}

	// 准备参数
	params := make([]interface{}, len(args))
	for i, arg := range args {
		switch arg.Kind {
		case WasmValueI32:
			params[i] = arg.Data.(int32)
		case WasmValueI64:
			params[i] = arg.Data.(int64)
		case WasmValueF32:
			params[i] = arg.Data.(float32)
		case WasmValueF64:
			params[i] = arg.Data.(float64)
		case WasmValueAnyRef:
			params[i] = arg.Data.(int32) // 引用类型通常表示为 i32
		case WasmValueFuncRef:
			params[i] = arg.Data.(int32) // 函数引用表示为 i32
		case WasmValueExternRef:
			params[i] = arg.Data.(int32) // 外部引用表示为 i32
		case WasmValueString:
			str := arg.Data.(string)
			ptr, err := c.Instance.WriteString(str)
			if err != nil {
				return nil, fmt.Errorf("failed to write string to memory: %v", err)
			}
			params[i] = ptr
		default:
			return nil, fmt.Errorf("unsupported argument type: %v", arg.Kind)
		}
	}

	// 准备结果数组
	results := make([]interface{}, 16) // 支持最多16个返回值

	// 调用函数
	err := c.Instance.CallFuncV(funcName, uint32(len(results)), results, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to call function %s: %v", funcName, err)
	}

	// 转换返回值
	wasmResults := make([]WasmValue, 0)
	for i, result := range results {
		if result == nil {
			break
		}

		fmt.Printf("result[%d] type: %T, value: %v\n", i, result, result)
		switch v := result.(type) {
		case int32:
			if funcName == "handle" {
				str, err := c.Instance.ReadString(result.(int32))
				if err != nil {
					return nil, fmt.Errorf("failed to read string from memory: %v", err)
				}
				defer c.Instance.FreeString(result.(int32))
				wasmResults = append(wasmResults, WasmValue{Kind: WasmValueString, Data: str})
			} else {
				wasmResults = append(wasmResults, WasmValue{Kind: WasmValueI32, Data: result})
			}
		case int64:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueI64, Data: result})
		case float32:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueF32, Data: result})
		case float64:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueF64, Data: result})
		case string:
			str, err := c.Instance.ReadString(result.(int32))
			if err != nil {
				return nil, fmt.Errorf("failed to read string from memory: %v", err)
			}
			defer c.Instance.FreeString(result.(int32))
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueString, Data: str})
		case unsafe.Pointer:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueAnyRef, Data: int32(uintptr(result.(unsafe.Pointer)))})
		case *int32: // funcref
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueFuncRef, Data: int32(uintptr(unsafe.Pointer(result.(*int32))))})
		case *interface{}: // externref
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueExternRef, Data: int32(uintptr(unsafe.Pointer(result.(*interface{}))))})
		default:
			return nil, fmt.Errorf("unsupported return type: %T", v)
		}
	}

	return wasmResults, nil
}
