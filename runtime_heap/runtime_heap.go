package runtime_heap

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_error"
)

var TYPE_STRING = "STRING"
var TYPE_NUMBER = "NUMBER"
var TYPE_FUNCTION = "FUNCTION"
var TYPE_NATIVE_FUNCTION = "NATIVE_FUNCTION"
var TYPE_UNKNOWN = "UNKNOWN"

type VariableValue struct {
	ValueType           string
	StringValue         string
	NumberValue         float64
	FunctionValue       *ast_node.ASTNode
	FunctionClosureHeap *Heap
	NativeFunctionName  string
}

type Heap struct {
	parentHeap *Heap
	values     map[string]*VariableValue
}

func CreateHeap() Heap {
	var values = map[string]*VariableValue{}

	return Heap{
		parentHeap: nil,
		values:     values,
	}
}

func (heap *Heap) CreateVariable(name string) error {
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
	prevVariable.NativeFunctionName = variable.NativeFunctionName
	prevVariable.FunctionClosureHeap = variable.FunctionClosureHeap

	return nil
}

func (heap *Heap) GetVariable(name string) *VariableValue {
	var variable = heap.values[name]

	if variable == nil && heap.parentHeap != nil {
		return heap.parentHeap.GetVariable(name)
	}

	return variable
}

func (heap *Heap) SetParentHeap(parent interface{}) {
	heap.parentHeap = parent.(*Heap)
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
			StringValue: strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", variable.NumberValue), "0"), "."),
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
