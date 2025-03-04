package main

import (
	"aolite-wamr-evn/core"
	"fmt"
	"log"
	"os"
)

func main() {
	// Create context
	ctx, err := core.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	// Load WASM file
	wasmBytes, err := os.ReadFile("wasm/process.wasm")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize runtime
	err = ctx.InitRuntime(wasmBytes)
	if err != nil {
		log.Fatal(err)
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
	// args := []core.WasmValue{
	// 	{Kind: core.WasmValueF64, Data: float64(5)},
	// }

	// results, err := ctx.CallFunction("fac", args)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Result: %v\n", results[0].Data)
}
