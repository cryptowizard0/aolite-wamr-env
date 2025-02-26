package main

import (
	"fmt"
	"os"

	"github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go/wamr"
)

func main() {
	var module *wamr.Module
	var instance *wamr.Instance
	var err error

	err = wamr.Runtime().Init()
	if err != nil {
		fmt.Println("Failed to init runtime:", err)
		return
	}
	defer wamr.Runtime().Destroy()

	// 设置日志级别
	wamr.Runtime().SetLogLevel(wamr.LOG_LEVEL_WARNING)

	// 读取 wasm 文件
	fmt.Print("Loading wasm module\n")
	byteCode, err := os.ReadFile("wasm/process.wasm")
	if err != nil {
		fmt.Println("Failed to read wasm file:", err)
		return
	}

	// 设置更详细的日志级别以便调试
	// wamr.Runtime().SetLogLevel(wamr.LOG_LEVEL_DEBUG)

	module, err = wamr.NewModule(byteCode)
	if err != nil {
		fmt.Println("Failed to load module:", err)
		return
	}
	defer module.Destroy()

	// 实例化 wasm 模块
	fmt.Print("Instantiating wasm module\n")
	instance, err = wamr.NewInstance(module, 83886080, 17179869184)
	if err != nil {
		fmt.Println("Failed to instantiate module:", err)
		return
	}
	defer instance.Destroy()
	fmt.Println("ok")

}
