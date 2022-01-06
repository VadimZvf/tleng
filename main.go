package main

import (
	"fmt"
	"os"

	"github.com/VadimZvf/golang/parser"
	"github.com/VadimZvf/golang/runtime"
	"github.com/VadimZvf/golang/runtime_bridge_cli"
	"github.com/VadimZvf/golang/source_file"
)

func main() {
	filePath := os.Args[1]

	var src = source_file.GetSource(filePath)
	var parser = parser.CreateParser(src)

	var astRoot, astError = parser.Parse(true)

	if astError != nil {
		fmt.Println(astError)
	}

	if astError != nil {
		os.Exit(1)
	}

	var bridge = runtime_bridge_cli.CreateBridge()
	var rt = runtime.CreateRuntime(&bridge, true)
	var runtimeErr = rt.Run(astRoot)

	if runtimeErr != nil {
		fmt.Println(runtimeErr)
	}
}
