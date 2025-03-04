package core

import (
	"errors"
	"fmt"

	"github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go/wamr"
)

// NewContext 创建新的WASM上下文
func NewContext() (*Context, error) {
	runtime := wamr.Runtime()
	if err := runtime.Init(); err != nil {
		return nil, fmt.Errorf("failed to init runtime: %v", err)
	}

	// 设置日志级别
	runtime.SetLogLevel(wamr.LOG_LEVEL_DEBUG)

	return &Context{
		Initialized: false,
	}, nil
}

// Close 清理并释放WASM上下文
func (c *Context) Close() {
	if c.Instance != nil {
		c.Instance.Destroy()
	}
	if c.Module != nil {
		c.Module.Destroy()
	}
	wamr.Runtime().Destroy()

}

// InitRuntime 初始化WASM运行时
func (c *Context) InitRuntime(wasmBytes []byte) error {
	if len(wasmBytes) == 0 {
		return errors.New("wasm bytes is empty")
	}

	// 创建模块
	module, err := wamr.NewModule(wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to create module: %v", err)
	}
	c.Module = module

	// 创建实例 (设置内存限制)
	instance, err := wamr.NewInstance(module, 83886080, 17179869184) // ~80MB, ~16GB
	if err != nil {
		return fmt.Errorf("failed to create instance: %v", err)
	}
	c.Instance = instance

	c.Initialized = true
	return nil
}
