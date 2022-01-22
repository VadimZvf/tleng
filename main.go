package main

import (
	"syscall/js"

	"github.com/VadimZvf/golang/parser"
	"github.com/VadimZvf/golang/parser_error_printer"
	"github.com/VadimZvf/golang/runtime"
	"github.com/VadimZvf/golang/runtime_bridge_web"
	"github.com/VadimZvf/golang/runtime_error_printer"
	"github.com/VadimZvf/golang/source_string"
	"github.com/VadimZvf/golang/stdout_web"
)

type iStdout interface {
	Print(line string)
	PrintError(line string)
}

func main() {
	// if len(os.Args) > 0 {
	// 	filePath := os.Args[1]

	// 	TlengRunFile(filePath)
	// }

	js.Global().Set("TlengRun", js.FuncOf(TlengWebRun))

	<-make(chan bool)
}

func TlengWebRun(this js.Value, args []js.Value) interface{} {
	codeText := args[0].String() // get the parameters
	var src = source_string.GetSource(codeText)
	var bridge = runtime_bridge_web.CreateBridge()
	var stdout = stdout_web.CreateStdoutWeb()
	Run(src, &stdout, &bridge)

	return nil
}

// func TlengRunFile(filePath string) interface{} {
// 	var src = source_file.GetSource(filePath)
// 	var bridge = runtime_bridge_cli.CreateBridge()
// 	var stdout = stdout.CreateStdout()
// 	Run(src, &stdout, &bridge)

// 	return nil
// }

func Run(source parser.ISource, stdout iStdout, bridge runtime.IBridge) interface{} {
	var parser = parser.CreateParser(source, stdout)

	var astRoot, astError = parser.Parse(false)

	if astError != nil {
		parser_error_printer.PrintError(parser.GetSourceCode(), stdout, astError)
		return nil
	}

	var rt = runtime.CreateRuntime(bridge)
	var runtimeErr = rt.Run(astRoot)

	if runtimeErr != nil {
		runtime_error_printer.PrintError(parser.GetSourceCode(), stdout, runtimeErr)
	}

	return nil
}
