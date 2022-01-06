package runtime

import (
	"fmt"
	"strconv"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_error"
	"github.com/VadimZvf/golang/runtime_heap"
)

type iHeap interface {
	CreateVariable(name string) error
	SetVariableValue(name string, variableType string, value string) error
	GetVariable(name string) (*runtime_heap.VariableValue, error)
}

type iBridge interface {
	Print(value string)
}

type Runtime struct {
	heap   iHeap
	bridge iBridge
}

func CreateRuntime(bridge iBridge, isDebug bool) Runtime {
	var heap = runtime_heap.CreateHeap(isDebug)

	return Runtime{
		heap:   &heap,
		bridge: bridge,
	}
}

type Visitor func(*ast_node.ASTNode, iHeap) (*runtime_heap.VariableValue, error)

func (runtime *Runtime) Run(ast *ast_node.ASTNode) error {
	var _, err = visitNode(ast, runtime.heap)

	return err
}

func visitNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var visitors = map[string]Visitor{
		ast_node.AST_NODE_CODE_ROOT:                     visitRootNode,
		ast_node.AST_NODE_CODE_VARIABLE_DECLARATION:     visitVariableDeclarationNode,
		ast_node.AST_NODE_CODE_ASSIGNMENT:               visitAssignmentNode,
		ast_node.AST_NODE_CODE_REFERENCE:                visitReferenceNode,
		ast_node.AST_NODE_CODE_NUMBER:                   visitNumberNode,
		ast_node.AST_NODE_CODE_STRING:                   visitStringNode,
		ast_node.AST_NODE_CODE_BINARY_EXPRESSION:        visitBinaryExpressionNode,
		ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION: visitParenthesizedExpressionNode,
	}

	var visitor = visitors[node.Code]

	if visitor == nil {
		return nil, runtime_error.RuntimeError{
			Message: "Unknown ast node: " + node.Code,
		}
	}

	return visitor(node, heap)
}

func visitRootNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	for _, bodyNode := range node.Body {
		var _, bodyNodeErr = visitNode(bodyNode, heap)

		if bodyNodeErr != nil {
			return nil, bodyNodeErr
		}
	}

	return nil, nil
}

func visitVariableDeclarationNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var variableNameParam = ast_node.GetVariableNameParam(node)

	if variableNameParam == nil {
		return nil, runtime_error.CreateError(
			"Cannot define variable without name",
			node,
		)
	}

	var err = heap.CreateVariable(variableNameParam.Value)

	return nil, err
}

func visitAssignmentNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var variableReferenceNode = node.Body[0]

	if variableReferenceNode == nil {
		return nil, runtime_error.CreateError(
			"Cannot set variable value without name",
			node,
		)
	}

	var variableValueNode = node.Body[1]

	if variableValueNode == nil {
		return nil, runtime_error.CreateError(
			"Variable value not defined",
			node,
		)
	}

	var variableName, getVariableNameErr = getVariableName(variableReferenceNode)

	if getVariableNameErr != nil {
		return nil, getVariableNameErr
	}

	var value, variableValueError = visitNode(variableValueNode, heap)

	if variableValueError != nil {
		return nil, variableValueError
	}

	if value == nil {
		return nil, runtime_error.CreateError(
			"Cannot get value from right node",
			variableValueNode,
		)
	}

	var err = heap.SetVariableValue(variableName, value.ValueType, value.Value)

	return nil, err
}

func visitReferenceNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var variableName, variableNameErr = getVariableName(node)

	if variableNameErr != nil {
		return nil, variableNameErr
	}

	var value, err = heap.GetVariable(variableName)

	return value, err
}

func visitNumberNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var numberValue = ast_node.GetNumberValueParam(node)

	if numberValue == nil {
		return nil, runtime_error.CreateError(
			"Cannot get number value",
			node,
		)
	}

	return &runtime_heap.VariableValue{Value: numberValue.Value, ValueType: runtime_heap.TYPE_NUMBER}, nil
}

func visitStringNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var stringValue = ast_node.GetStringValueParam(node)

	if stringValue == nil {
		return nil, runtime_error.CreateError(
			"Cannot get number value",
			node,
		)
	}

	return &runtime_heap.VariableValue{Value: stringValue.Value, ValueType: runtime_heap.TYPE_STRING}, nil
}

func visitParenthesizedExpressionNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var bodyNode = node.Body[0]

	if bodyNode == nil {
		return nil, runtime_error.CreateError(
			"Cannot get body node of parenthesized expression",
			node,
		)
	}

	return visitNode(bodyNode, heap)
}

func visitBinaryExpressionNode(node *ast_node.ASTNode, heap iHeap) (*runtime_heap.VariableValue, error) {
	var leftNode = node.Body[0]

	if leftNode == nil {
		return nil, runtime_error.CreateError(
			"Cannot get left node of binary expression",
			node,
		)
	}

	var rightNode = node.Body[1]

	if rightNode == nil {
		return nil, runtime_error.CreateError(
			"Cannot get right node of binary expression",
			node,
		)
	}

	var leftNodeValue, leftNodeError = visitNode(leftNode, heap)

	if leftNodeError != nil {
		return nil, leftNodeError
	}

	var rightNodeValue, rightNodeError = visitNode(rightNode, heap)

	if rightNodeError != nil {
		return nil, rightNodeError
	}

	var expressionType = ast_node.GetBinaryExpressionTypeParam(node)

	if expressionType == nil {
		return nil, runtime_error.CreateError(
			"Binary expression type not defined",
			node,
		)
	}

	if leftNodeValue.ValueType == runtime_heap.TYPE_NUMBER && rightNodeValue.ValueType == runtime_heap.TYPE_NUMBER {
		var leftNumberValue, leftNumberParsError = strconv.ParseFloat(leftNodeValue.Value, 32)
		if leftNumberParsError != nil {
			return nil, runtime_error.CreateError(
				"Failed parse number value at left node",
				leftNode,
			)
		}

		var rightNumberValue, rightNumberParsError = strconv.ParseFloat(rightNodeValue.Value, 32)
		if rightNumberParsError != nil {
			return nil, runtime_error.CreateError(
				"Failed parse number value at right node",
				rightNode,
			)
		}

		switch expressionType.Value {
		case "+":
			return &runtime_heap.VariableValue{
				Value:     fmt.Sprintf("%f", leftNumberValue+rightNumberValue),
				ValueType: runtime_heap.TYPE_NUMBER,
			}, nil
		case "-":
			return &runtime_heap.VariableValue{
				Value:     fmt.Sprintf("%f", leftNumberValue-rightNumberValue),
				ValueType: runtime_heap.TYPE_NUMBER,
			}, nil
		case "/":
			return &runtime_heap.VariableValue{
				Value:     fmt.Sprintf("%f", leftNumberValue/rightNumberValue),
				ValueType: runtime_heap.TYPE_NUMBER,
			}, nil
		case "*":
			return &runtime_heap.VariableValue{
				Value:     fmt.Sprintf("%f", leftNumberValue*rightNumberValue),
				ValueType: runtime_heap.TYPE_NUMBER,
			}, nil
		}
	}

	switch expressionType.Value {
	case "+":
		return &runtime_heap.VariableValue{
			Value:     leftNodeValue.Value + rightNodeValue.Value,
			ValueType: runtime_heap.TYPE_STRING,
		}, nil
	}

	return nil, runtime_error.CreateError(
		"Unknown binary expression. Received: "+expressionType.Value,
		node,
	)
}

func getVariableName(node *ast_node.ASTNode) (string, error) {
	if node.Code == ast_node.AST_NODE_CODE_REFERENCE {
		var variableNameParam = ast_node.GetVariableNameParam(node)

		if variableNameParam == nil {
			return "", runtime_error.CreateError(
				"Reference without variable without name",
				node,
			)
		}

		return variableNameParam.Value, nil
	}

	return "", runtime_error.CreateError(
		"Cannot get name from node: "+node.Code,
		node,
	)
}
