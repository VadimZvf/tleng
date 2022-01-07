package runtime

import (
	"strconv"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/runtime_error"
	"github.com/VadimZvf/golang/runtime_heap"
)

type iHeap interface {
	CreateVariable(name string) error
	SetVariable(name string, variable *runtime_heap.VariableValue) error
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
	var heap = runtime_heap.CreateHeap(isDebug, bridge)

	return Runtime{
		heap:   &heap,
		bridge: bridge,
	}
}

type Visitor func(*ast_node.ASTNode) (*runtime_heap.VariableValue, error)

func (runtime *Runtime) Run(ast *ast_node.ASTNode) error {
	var _, err = runtime.visitNode(ast)

	return err
}

func (runtime *Runtime) visitNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var visitors = map[string]Visitor{
		ast_node.AST_NODE_CODE_ROOT:                     runtime.visitRootNode,
		ast_node.AST_NODE_CODE_VARIABLE_DECLARATION:     runtime.visitVariableDeclarationNode,
		ast_node.AST_NODE_CODE_ASSIGNMENT:               runtime.visitAssignmentNode,
		ast_node.AST_NODE_CODE_REFERENCE:                runtime.visitReferenceNode,
		ast_node.AST_NODE_CODE_NUMBER:                   runtime.visitNumberNode,
		ast_node.AST_NODE_CODE_STRING:                   runtime.visitStringNode,
		ast_node.AST_NODE_CODE_BINARY_EXPRESSION:        runtime.visitBinaryExpressionNode,
		ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION: runtime.visitParenthesizedExpressionNode,
	}

	var visitor = visitors[node.Code]

	if visitor == nil {
		return nil, runtime_error.CreateError("Unknown ast node: "+node.Code,
			node,
		)
	}

	return visitor(node)
}

func (runtime *Runtime) visitRootNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	for _, bodyNode := range node.Body {
		var _, bodyNodeErr = runtime.visitNode(bodyNode)

		if bodyNodeErr != nil {
			return nil, bodyNodeErr
		}
	}

	return nil, nil
}

func (runtime *Runtime) visitVariableDeclarationNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var variableNameParam = ast_node.GetVariableNameParam(node)

	if variableNameParam == nil {
		return nil, runtime_error.CreateError(
			"Cannot define variable without name",
			node,
		)
	}

	var err = runtime.heap.CreateVariable(variableNameParam.Value)

	return nil, err
}

func (runtime *Runtime) visitAssignmentNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
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

	var value, variableValueError = runtime.visitNode(variableValueNode)

	if variableValueError != nil {
		return nil, variableValueError
	}

	if value == nil {
		return nil, runtime_error.CreateError(
			"Cannot get value from right node",
			variableValueNode,
		)
	}

	var err = runtime.heap.SetVariable(variableName, value)

	return nil, err
}

func (runtime *Runtime) visitReferenceNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var variableName, variableNameErr = getVariableName(node)

	if variableNameErr != nil {
		return nil, variableNameErr
	}

	var value, err = runtime.heap.GetVariable(variableName)

	return value, err
}

func (runtime *Runtime) visitNumberNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var numberValue = ast_node.GetNumberValueParam(node)

	if numberValue == nil {
		return nil, runtime_error.CreateError(
			"Cannot get number value",
			node,
		)
	}

	var number, numberParsError = strconv.ParseFloat(numberValue.Value, 64)

	if numberParsError != nil {
		return nil, runtime_error.CreateError(
			"Failed parse number value at left node",
			node,
		)
	}

	return &runtime_heap.VariableValue{NumberValue: number, ValueType: runtime_heap.TYPE_NUMBER}, nil
}

func (runtime *Runtime) visitStringNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var stringValue = ast_node.GetStringValueParam(node)

	if stringValue == nil {
		return nil, runtime_error.CreateError(
			"Cannot get number value",
			node,
		)
	}

	return &runtime_heap.VariableValue{StringValue: stringValue.Value, ValueType: runtime_heap.TYPE_STRING}, nil
}

func (runtime *Runtime) visitParenthesizedExpressionNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var bodyNode = node.Body[0]

	if bodyNode == nil {
		return nil, runtime_error.CreateError(
			"Cannot get body node of parenthesized expression",
			node,
		)
	}

	return runtime.visitNode(bodyNode)
}

func (runtime *Runtime) visitBinaryExpressionNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
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

	var leftNodeValue, leftNodeError = runtime.visitNode(leftNode)

	if leftNodeError != nil {
		return nil, leftNodeError
	}

	var rightNodeValue, rightNodeError = runtime.visitNode(rightNode)

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
		var leftNumberValue = leftNodeValue.NumberValue
		var rightNumberValue = rightNodeValue.NumberValue

		switch expressionType.Value {
		case "+":
			return &runtime_heap.VariableValue{
				NumberValue: leftNumberValue + rightNumberValue,
				ValueType:   runtime_heap.TYPE_NUMBER,
			}, nil
		case "-":
			return &runtime_heap.VariableValue{
				NumberValue: leftNumberValue - rightNumberValue,
				ValueType:   runtime_heap.TYPE_NUMBER,
			}, nil
		case "/":
			return &runtime_heap.VariableValue{
				NumberValue: leftNumberValue / rightNumberValue,
				ValueType:   runtime_heap.TYPE_NUMBER,
			}, nil
		case "*":
			return &runtime_heap.VariableValue{
				NumberValue: leftNumberValue * rightNumberValue,
				ValueType:   runtime_heap.TYPE_NUMBER,
			}, nil
		}
	}

	var leftString, leftNodeCastErr = runtime_heap.CastToString(leftNodeValue)

	if leftNodeCastErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot convert variable value to string. Received: "+leftNodeValue.ValueType,
			node,
		), leftNodeCastErr)
	}

	var rightString, rightNodeCastErr = runtime_heap.CastToString(rightNodeValue)

	if rightNodeCastErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot convert variable value to string. Received: "+rightString.ValueType,
			node,
		), rightNodeCastErr)
	}

	switch expressionType.Value {
	case "+":
		return &runtime_heap.VariableValue{
			StringValue: leftString.StringValue + rightString.StringValue,
			ValueType:   runtime_heap.TYPE_STRING,
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
