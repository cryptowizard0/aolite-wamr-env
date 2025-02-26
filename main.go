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

	// Set log level
	wamr.Runtime().SetLogLevel(wamr.LOG_LEVEL_WARNING)
	// Or set more detailed log level for debugging
	wamr.Runtime().SetLogLevel(wamr.LOG_LEVEL_DEBUG)

	// loading wasm file
	fmt.Print("Loading wasm module\n")
	byteCode, err := os.ReadFile("wasm/process.wasm")
	if err != nil {
		fmt.Println("Failed to read wasm file:", err)
		return
	}

	module, err = wamr.NewModule(byteCode)
	if err != nil {
		fmt.Println("Failed to load module:", err)
		return
	}
	defer module.Destroy()

	// Instantiating wasm module
	fmt.Print("Instantiating wasm module\n")
	instance, err = wamr.NewInstance(module, 83886080, 17179869184)
	if err != nil {
		fmt.Println("Failed to instantiate module:", err)
		return
	}
	defer instance.Destroy()
	fmt.Println("ok")

}
