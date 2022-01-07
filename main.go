package main

import (
	"fmt"
	"os"

	"github.com/VadimZvf/golang/parser"
	"github.com/VadimZvf/golang/runtime"
	"github.com/VadimZvf/golang/runtime_bridge_cli"
	"github.com/VadimZvf/golang/runtime_error_printer"
	"github.com/VadimZvf/golang/source_file"
	"github.com/VadimZvf/golang/stdout"
)

func main() {
	filePath := os.Args[1]

	var src = source_file.GetSource(filePath)
	var stdout = stdout.CreateStdout()
	var parser = parser.CreateParser(src, &stdout)

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
		runtime_error_printer.PrintError(parser.GetSourceCode(), &stdout, runtimeErr)
	}
}
