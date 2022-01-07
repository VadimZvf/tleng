package runtime_heap

import (
	"fmt"
	"strconv"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_error"
)

var TYPE_STRING = "STRING"
var TYPE_NUMBER = "NUMBER"
var TYPE_FUNCTION = "FUNCTION"
var TYPE_UNKNOWN = "UNKNOWN"

type VariableValue struct {
	ValueType     string
	StringValue   string
	NumberValue   float64
	FunctionValue *ast_node.ASTNode
}

type iBridge interface {
	Print(value string)
}

type Heap struct {
	idDebug bool

	values map[string]*VariableValue
	bridge iBridge
}

func CreateHeap(idDebug bool, bridge iBridge) Heap {
	var values = map[string]*VariableValue{}

	return Heap{
		idDebug,
		values,
		bridge,
	}
}

func (heap *Heap) CreateVariable(name string) error {
	if heap.idDebug {
		heap.bridge.Print("Create variable with name: " + name)
	}

	var prevVariable = heap.values[name]

	if prevVariable != nil {
		return runtime_error.RuntimeError{
			Message: "Variable already declared",
		}
	}

	heap.values[name] = &VariableValue{
		ValueType: TYPE_UNKNOWN,
	}

	return nil
}

func (heap *Heap) SetVariable(name string, variable *VariableValue) error {
	if heap.idDebug {
		heap.bridge.Print("Set variable value name: " + name)
		logVariable(variable, heap.bridge)
	}

	var prevVariable = heap.values[name]

	if prevVariable == nil {
		return runtime_error.RuntimeError{
			Message: "Variable not declared",
		}
	}

	prevVariable.ValueType = variable.ValueType
	prevVariable.NumberValue = variable.NumberValue
	prevVariable.StringValue = variable.StringValue
	prevVariable.FunctionValue = variable.FunctionValue

	return nil
}

func (heap *Heap) GetVariable(name string) (*VariableValue, error) {
	if heap.idDebug {
		heap.bridge.Print("Get variable value name: " + name)
	}

	var prevVariable = heap.values[name]

	if prevVariable == nil {
		return nil, runtime_error.RuntimeError{
			Message: "Variable not defined",
		}
	}

	return prevVariable, nil
}

func CastToNumber(variable *VariableValue) (*VariableValue, error) {
	if variable.ValueType == TYPE_NUMBER {
		return variable, nil
	}

	if variable.ValueType == TYPE_STRING {
		var number, numberParsError = strconv.ParseFloat(variable.StringValue, 64)

		if numberParsError != nil {
			return nil, numberParsError
		}

		return &VariableValue{
			ValueType:   TYPE_NUMBER,
			NumberValue: number,
		}, nil
	}

	if variable.ValueType == TYPE_UNKNOWN {
		return &VariableValue{
			ValueType:   TYPE_NUMBER,
			NumberValue: 0,
		}, nil
	}

	return nil, runtime_error.RuntimeError{
		Message: "Cannot cast variable to number. Type: " + variable.ValueType,
	}
}

func CastToString(variable *VariableValue) (*VariableValue, error) {
	if variable.ValueType == TYPE_STRING {
		return variable, nil
	}

	if variable.ValueType == TYPE_NUMBER {
		return &VariableValue{
			ValueType:   TYPE_STRING,
			StringValue: fmt.Sprintf("%f", variable.NumberValue),
		}, nil
	}

	if variable.ValueType == TYPE_UNKNOWN {
		return &VariableValue{
			ValueType:   TYPE_STRING,
			StringValue: "",
		}, nil
	}

	return nil, runtime_error.RuntimeError{
		Message: "Cannot cast variable to number. Type: " + variable.ValueType,
	}
}

func logVariable(variable *VariableValue, bridge iBridge) {
	switch variable.ValueType {
	case TYPE_STRING:
		bridge.Print(fmt.Sprintf("value type: %s, value: %s\n", variable.ValueType, variable.StringValue))
		return
	case TYPE_NUMBER:
		bridge.Print(fmt.Sprintf("value type: %s, value: %f\n", variable.ValueType, variable.NumberValue))
		return
	case TYPE_FUNCTION:
		var runctionNameParam = ast_node.GetFunctionNameParam(variable.FunctionValue)
		bridge.Print(fmt.Sprintf("value type: %s, name: %s\n", variable.ValueType, runctionNameParam.Value))
		return
	case TYPE_UNKNOWN:
		bridge.Print(fmt.Sprintf("value type: %s\n", variable.ValueType))
		return
	}

	bridge.Print(fmt.Sprintf("value type: %s\n", variable.ValueType))
}
