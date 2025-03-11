package main

import (
	"aolite-wamr-evn/core"
	"fmt"
	"os"
	"testing"

	"github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go/wamr"
)

func TestWasmFactorial(t *testing.T) {
	// Create context
	ctx, err := core.NewContext()
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}
	defer ctx.Close()

	// Load WASM file
	wasmBytes, err := os.ReadFile("wasm/test-64.wasm")
	if err != nil {
		t.Fatalf("Failed to read WASM file: %v", err)
	}

	// Initialize runtime
	err = ctx.InitRuntime(wasmBytes)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}

	// Call function
	args := []core.WasmValue{
		{
			Kind: core.WasmValueF64,
			Data: 5.0,
		},
	}

	results, err := ctx.CallFunction("fac", args)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	// Verify result
	expected := float64(120) // factorial of 5
	if results[0].Data != expected {
		t.Errorf("Expected %v, got %v", expected, results[0].Data)
	}
}

func TestWasmPow(t *testing.T) {
	// Create context
	ctx, err := core.NewContext()
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}
	defer ctx.Close()

	err = wamr.RegisterImportFunctions()
	if err != nil {
		t.Fatalf("Failed to register WASI functions: %v", err)
	}

	// Load WASM file
	wasmBytes, err := os.ReadFile("wasm/pow_calculator.wasm")
	if err != nil {
		t.Fatalf("Failed to read WASM file: %v", err)
	}

	// Initialize runtime
	err = ctx.InitRuntime(wasmBytes)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}

	// get exports count from wasm file
	count, _ := ctx.GetExportCount()
	fmt.Println("export count: ", count)

	// print export info from wasm file
	for i := int32(0); i < count; i++ {
		export, _ := ctx.GetExportType(i)
		fmt.Printf("Export #%d: name=%s, kind=%d\n",
			i, export.Name, export.Kind)
	}

	// Call function
	args := []core.WasmValue{
		{
			Kind: core.WasmValueI32,
			Data: int32(2),
		},
		{
			Kind: core.WasmValueI32,
			Data: int32(2),
		},
	}

	results, err := ctx.CallFunction("pow", args)
	if err != nil {
		t.Fatalf("Failed to call function: %v\n", err)
	}

	fmt.Printf("results: %#v \n", results) // 打印 results结构体

	// Verify result
	expected := int32(4) // factorial of 5
	if results[0].Data != expected {
		t.Fatalf("Expected %v, got %v \n", expected, results[0].Data)
	}
	fmt.Println("all ok!")
}

func TestListExportedFunctions(t *testing.T) {
	// Create context
	ctx, err := core.NewContext()
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}
	defer ctx.Close()

	// Load WASM file
	wasmBytes, err := os.ReadFile("wasm/test-64.wasm")
	if err != nil {
		t.Fatalf("Failed to read WASM file: %v", err)
	}

	// Initialize runtime
	err = ctx.InitRuntime(wasmBytes)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}

	// // Get exported functions
	count, _ := ctx.GetExportCount()
	fmt.Println("export count: ", count)

	// print export info from wasm file
	for i := int32(0); i < count; i++ {
		export, _ := ctx.GetExportType(i)
		fmt.Printf("Export #%d: name=%s, kind=%d\n",
			i, export.Name, export.Kind)
	}
}
