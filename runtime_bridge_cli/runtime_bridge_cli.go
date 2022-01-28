package runtime_bridge_cli

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_heap"
)

type Bridge struct {
}

func CreateBridge() Bridge {
	return Bridge{}
}

func (bridge *Bridge) Print(args ...*runtime_heap.VariableValue) {
	for _, arg := range args {
		printArg(arg)
	}
}

func printArg(variable *runtime_heap.VariableValue) {
	if variable.ValueType == runtime_heap.TYPE_STRING {
		fmt.Println(variable.StringValue)
	}

	if variable.ValueType == runtime_heap.TYPE_NUMBER {
		fmt.Println(variable.NumberValue)
	}

	if variable.ValueType == runtime_heap.TYPE_BOOLEAN {
		fmt.Println(variable.BooleanValue)
	}

	if variable.ValueType == runtime_heap.TYPE_FUNCTION {
		var functionName = ast_node.GetFunctionNameParam(variable.FunctionValue)
		fmt.Println("function " + functionName.Value)
	}

	if variable.ValueType == runtime_heap.TYPE_NATIVE_FUNCTION {
		fmt.Println("native code")
	}

	if variable.ValueType == runtime_heap.TYPE_UNKNOWN {
		fmt.Println("unknown")
	}
}
