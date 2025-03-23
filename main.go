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
	//print_string()
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
	// count, _ := ctx.GetExportCount()
	// fmt.Println("export count: ", count)

	// // print export info from wasm file
	// for i := int32(0); i < count; i++ {
	// 	export, _ := ctx.GetExportType(i)
	// 	fmt.Printf("Export #%d: name=%s, kind=%d\n",
	// 		i, export.Name, export.Kind)
	// }

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
        "Action": "Eval",
        "Block-Height": "100000",
        "Data": "print(\"hello eval\")",
        "Data-Protocol": "aolite",
        "From": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ",
        "Id": "l4Ya5H0RAHZIUoKZzd37AyAAXR7Vx30Hmf2CuOuoPzc",
        "Module": "0x84534",
        "Owner": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ",
        "Tags": [
            {
                "name": "From",
                "value": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ"
            },
            {
                "name": "Data-Protocol",
                "value": "aolite"
            },
            {
                "name": "Action",
                "value": "Eval"
            },
            {
                "name": "Target",
                "value": "2cnUXq0CXgMnjvpZn4w_2HW6zdHPKk6M1jAPBD5do20"
            },
            {
                "name": "Module",
                "value": "0x84534"
            },
            {
                "name": "Block-Height",
                "value": "100000"
            },
            {
                "name": "Data",
                "value": "print(\"hello eval\")"
            },
            {
                "name": "Variant",
                "value": "aolite.TN.1"
            },
            {
                "name": "Type",
                "value": "Message"
            },
            {
                "name": "Owner",
                "value": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ"
            },
            {
                "name": "Id",
                "value": "l4Ya5H0RAHZIUoKZzd37AyAAXR7Vx30Hmf2CuOuoPzc"
            }
        ],
        "Target": "2cnUXq0CXgMnjvpZn4w_2HW6zdHPKk6M1jAPBD5do20",
        "Type": "Message",
        "Variant": "aolite.TN.1"
    }`

	env := `{
		"Module": {
			"Owner": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ",
			"Id": "0x456",
			"Tags": [
			{
				"name": "Data-Protocol",
				"value": "aolite"
			},
			{
				"name": "Variant",
				"value": "aolite.TN.1"
			},
			{
				"name": "Type",
				"value": "Process"
			},
			{
				"name": "Module-Format",
				"value": "lua5.3"
			},
			{
				"name": "Memory-Limit",
				"value": "500-mb"
			},
			{
				"name": "Compute-Limit",
				"value": "9000000000000"
			},
			{
				"name": "Content-Type",
				"value": "text/plain"
			}
			]
		},
		"Process": {
			"Owner": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ",
			"Id": "l4Ya5H0RAHZIUoKZzd37AyAAXR7Vx30Hmf2CuOuoPzc",
			"Tags": [
			{
				"name": "Data-Protocol",
				"value": "aolite"
			},
			{
				"name": "Variant",
				"value": "aolite.TN.1"
			},
			{
				"name": "Type",
				"value": "Process"
			},
			{
				"name": "Module",
				"value": "LSjhdzBjyWuyUPe-g6PUzt8t1PUlw2FZ9SM3_hCh2Is"
			},
			{
				"name": "Scheduler",
				"value": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ"
			},
			{
				"name": "Authority",
				"value": "eIgnDk4vSKPe0lYB6yhCHDV1dOw3JgYHGocfj7WGrjQ"
			},
			{
				"name": "App-Name",
				"value": "aolite"
			},
			{
				"name": "Name",
				"value": "default"
			},
			{
				"name": "Content-Type",
				"value": "text/plain"
			},
			{
				"name": "Reference",
				"value": "0"
			}
			]
		}
	}`
	// env = "{}"
	// msg = "{}"
	handleArgs := []core.WasmValue{
		{Kind: core.WasmValueString, Data: msg},
		{Kind: core.WasmValueString, Data: env},
	}

	results, err := ctx.CallFunction("handle", handleArgs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Result: %v\n", results[0].Data)
}

func print_string() {
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
	wasmBytes, err := os.ReadFile("wasm/print_string-emcc.wasm")
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

	// Call function
	// 准备参数
	testString := "Hello from Go!"
	args := []core.WasmValue{
		{Kind: core.WasmValueString, Data: testString},
	}

	_, err = ctx.CallFunction("print_string", args)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Function called successfully!")
}
