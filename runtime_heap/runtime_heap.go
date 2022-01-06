package runtime_heap

import (
	"fmt"

	"github.com/VadimZvf/golang/runtime_error"
)

var TYPE_STRING = "STRING"
var TYPE_NUMBER = "NUMBER"
var TYPE_UNKNOWN = "UNKNOWN"

type VariableValue struct {
	ValueType string
	Value     string
}

type Heap struct {
	idDebug bool

	values map[string]*VariableValue
}

func CreateHeap(idDebug bool) Heap {
	var values = map[string]*VariableValue{}

	return Heap{
		idDebug,
		values,
	}
}

func (heap *Heap) CreateVariable(name string) error {
	if heap.idDebug {
		fmt.Printf("Create variable with name: %s\n", name)
	}

	var prevVariable = heap.values[name]

	if prevVariable != nil {
		return runtime_error.RuntimeError{
			Message: "Variable already declared",
		}
	}

	heap.values[name] = &VariableValue{
		Value:     "",
		ValueType: TYPE_UNKNOWN,
	}

	return nil
}

func (heap *Heap) SetVariableValue(name string, variableType string, value string) error {
	if heap.idDebug {
		fmt.Printf("Set variable value name: %s value: %s\n", name, value)
	}

	var prevVariable = heap.values[name]

	if prevVariable == nil {
		return runtime_error.RuntimeError{
			Message: "Variable not defined declared",
		}
	}

	prevVariable.ValueType = variableType
	prevVariable.Value = value

	return nil
}

func (heap *Heap) GetVariable(name string) (*VariableValue, error) {
	if heap.idDebug {
		fmt.Printf("Get variable value name: %s\n", name)
	}

	var prevVariable = heap.values[name]

	if prevVariable == nil {
		return nil, runtime_error.RuntimeError{
			Message: "Variable not defined declared",
		}
	}

	return prevVariable, nil
}
