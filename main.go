package main

import (
	"aolite-wamr-evn/core"
	"fmt"
	"log"
	"os"

	"github.com/bytecodealliance/wasm-micro-runtime/language-bindings/go/wamr"
)

func main() {
	LoadAO()
}

func LoadAO() {
	// Create context
	ctx, err := core.NewContext()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ctx.Close()

	// 在加载 WASM 之前注册 WASI 函数
	err = wamr.RegisterImportFunctions()
	if err != nil {
		log.Fatal("Failed to register WASI functions:", err)
		return
	}

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

	// get exports count from wasm file
	count, _ := ctx.GetExportCount()
	fmt.Println("export count: ", count)

	// print export info from wasm file
	for i := int32(0); i < count; i++ {
		export, _ := ctx.GetExportType(i)
		fmt.Printf("Export #%d: name=%s, kind=%d\n",
			i, export.Name, export.Kind)
	}

	// call main
	// args := []uint32{0, 0, 0}
	// err = ctx.Instance.CallFunc("main", 3, args)
	// if err != nil {
	// 	fmt.Println("call main error: ", err)
	// 	return
	// }
	// fmt.Println("call main ok!")
	// fmt.Printf("main function returns: args[0]=%d, args[1]=%d, args[2]=%d\n",
	// 	args[0], args[1], args[2])

	// Call function
	msg := `{
		"From": "FOOBAR",
		"To": "AOS",
		"Type": "Command",
		"Id": "MSG_ID",
		"Time": "2021-01-01T00:00:00Z",
		"Data": "return 1 + 1",
		"Metadata": {
			"Source": "test",
			"Target": "aos"
		}
	}`

	env := `{
		"Process": {
			"Id": "AOS",
			"Type": "Service",
			"Name": "Autonomous Operating System",
			"Version": "2.0.0"
		},
		"Runtime": {
			"Name": "WASM",
			"Version": "1.0.0",
			"Platform": "test"
		},
		"Context": {
			"RequestId": "REQ_ID",
			"SessionId": "SESSION_ID",
			"TraceId": "TRACE_ID"
		}
	}`
	handleArgs := []core.WasmValue{
		{Kind: core.WasmValueString, Data: msg},
		{Kind: core.WasmValueString, Data: env},
	}

	results, err := ctx.CallFunction("handle", handleArgs)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Result: %v\n", results[0].Data)
}
