package runtime_bridge_mock

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_heap"
)

type Bridge struct {
	log []string
}

func CreateBridge() Bridge {
	return Bridge{}
}

func (bridge *Bridge) Print(args ...*runtime_heap.VariableValue) {
	for _, arg := range args {
		bridge.saveLogArg(arg)
	}
}

func (bridge *Bridge) GetLastPring() string {
	if len(bridge.log) > 0 {
		return bridge.log[len(bridge.log)-1]
	}

	return ""
}

func (bridge *Bridge) saveLogArg(variable *runtime_heap.VariableValue) {
	if variable.ValueType == runtime_heap.TYPE_STRING {
		bridge.log = append(bridge.log, variable.StringValue)
	}

	if variable.ValueType == runtime_heap.TYPE_NUMBER {
		bridge.log = append(bridge.log, fmt.Sprintf("%f", variable.NumberValue))
	}

	if variable.ValueType == runtime_heap.TYPE_FUNCTION {
		var functionName = ast_node.GetFunctionNameParam(variable.FunctionValue)
		bridge.log = append(bridge.log, "function "+functionName.Value)
	}

	if variable.ValueType == runtime_heap.TYPE_NATIVE_FUNCTION {
		bridge.log = append(bridge.log, "native code")
	}

	if variable.ValueType == runtime_heap.TYPE_UNKNOWN {
		bridge.log = append(bridge.log, "unknown")
	}
}
