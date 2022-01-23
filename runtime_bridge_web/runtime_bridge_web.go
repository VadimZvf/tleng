package runtime_bridge_web

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_heap"
)

type JSPrint func(value string)

type Bridge struct {
	JSPrint JSPrint
}

func CreateBridge() Bridge {
	WindowJS := js.Global().Get("window")

	return Bridge{
		JSPrint: func(value string) {
			WindowJS.Call("TlengPrint", value)
		},
	}
}

func (bridge *Bridge) Print(args ...*runtime_heap.VariableValue) {
	for _, arg := range args {
		printArg(arg, bridge)
	}
}

func printArg(variable *runtime_heap.VariableValue, bridge *Bridge) {
	if variable.ValueType == runtime_heap.TYPE_STRING {
		bridge.JSPrint(variable.StringValue)
	}

	if variable.ValueType == runtime_heap.TYPE_NUMBER {
		bridge.JSPrint(strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", variable.NumberValue), "0"), "."))
	}

	if variable.ValueType == runtime_heap.TYPE_FUNCTION {
		var functionName = ast_node.GetFunctionNameParam(variable.FunctionValue)
		bridge.JSPrint("function " + functionName.Value)
	}

	if variable.ValueType == runtime_heap.TYPE_NATIVE_FUNCTION {
		bridge.JSPrint("native code")
	}

	if variable.ValueType == runtime_heap.TYPE_UNKNOWN {
		bridge.JSPrint("unknown")
	}
}
