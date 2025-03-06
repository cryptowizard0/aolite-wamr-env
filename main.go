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
		return
	}
	defer ctx.Close()

	// Load WASM file
	wasmBytes, err := os.ReadFile("wasm/process.wasm")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initialize runtime
	err = ctx.InitRuntime(wasmBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

	// err = ctx.RegisterWASIFunctions()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

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
		{Kind: core.WasmValueString, Data: ""},
		{Kind: core.WasmValueString, Data: ""},
	}

	results, err := ctx.CallFunction("handle", args)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Result: %v\n", results[0].Data)

	// args := []uint32{0, 0, 0}
	// err = ctx.Instance.CallFunc("main", 3, args)
	// if err != nil {
	// 	fmt.Println("call main error: ", err)
	// }

}
