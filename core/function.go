package core

import (
	"errors"
	"fmt"
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
	for _, result := range results {
		if result == nil {
			break
		}

		switch v := result.(type) {
		case int32:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueI32, Data: v})
		case int64:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueI64, Data: v})
		case float32:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueF32, Data: v})
		case float64:
			wasmResults = append(wasmResults, WasmValue{Kind: WasmValueF64, Data: v})
		default:
			return nil, fmt.Errorf("unsupported return type: %T", v)
		}
	}

	return wasmResults, nil
}
