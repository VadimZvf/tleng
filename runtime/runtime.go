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
	Print(args ...*runtime_heap.VariableValue)
}

type Runtime struct {
	heap   iHeap
	bridge iBridge
}

func CreateRuntime(bridge iBridge) Runtime {
	var heap = runtime_heap.CreateHeap()

	var rt = Runtime{
		heap:   &heap,
		bridge: bridge,
	}
	rt.defineEnvByBridge()

	return rt
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
		ast_node.AST_NODE_CODE_FUNCTION:                 runtime.visitFunctionNode,
		ast_node.AST_NODE_CODE_BINARY_EXPRESSION:        runtime.visitBinaryExpressionNode,
		ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION: runtime.visitParenthesizedExpressionNode,
		ast_node.AST_NODE_CODE_CALL_EXPRESSION:          runtime.visitCallExpressionNode,
		ast_node.AST_NODE_CODE_RETURN:                   runtime.visitReturnNode,
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
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get variable name for assertion",
			variableReferenceNode,
		), getVariableNameErr)
	}

	var value, variableValueError = runtime.visitNode(variableValueNode)

	if variableValueError != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get variable for assertion",
			variableValueNode,
		), variableValueError)
	}

	if value == nil {
		return nil, runtime_error.CreateError(
			"Cannot get value from right node",
			variableValueNode,
		)
	}

	var setVariableError = runtime.heap.SetVariable(variableName, value)

	if setVariableError != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot set variable",
			variableReferenceNode,
		), setVariableError)
	}

	return value, nil
}

func (runtime *Runtime) visitReferenceNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var variableName, variableNameErr = getVariableName(node)

	if variableNameErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get variable name",
			node,
		), variableNameErr)
	}

	var value, getVariableErr = runtime.heap.GetVariable(variableName)

	if getVariableErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get variable reference",
			node,
		), getVariableErr)
	}

	return value, nil
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

func (runtime *Runtime) visitReturnNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var bodyNode = node.Body[0]

	if bodyNode != nil {
		return runtime.visitNode(bodyNode)
	}

	return nil, nil
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
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get left node value",
			leftNode,
		), leftNodeError)
	}

	var rightNodeValue, rightNodeError = runtime.visitNode(rightNode)

	if rightNodeError != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get right node value",
			rightNode,
		), rightNodeError)

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

func (runtime *Runtime) visitFunctionNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var functionNameParam = ast_node.GetFunctionNameParam(node)

	if functionNameParam == nil {
		return nil, runtime_error.CreateError(
			"Cannot define function without name",
			node,
		)
	}

	var createVariableForFuncErr = runtime.heap.CreateVariable(functionNameParam.Value)

	if createVariableForFuncErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot define function with name: "+functionNameParam.Value,
			node,
		), createVariableForFuncErr)
	}

	var functionVariable, getFuncVariableErr = runtime.heap.GetVariable(functionNameParam.Value)

	if getFuncVariableErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get allocated variable for function with name: "+functionNameParam.Value,
			node,
		), getFuncVariableErr)
	}

	functionVariable.ValueType = runtime_heap.TYPE_FUNCTION
	functionVariable.FunctionValue = node

	var setFunctionError = runtime.heap.SetVariable(functionNameParam.Value, functionVariable)

	if setFunctionError != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot set function into variable with name: "+functionNameParam.Value,
			node,
		), setFunctionError)
	}

	return functionVariable, nil
}

func (runtime *Runtime) visitCallExpressionNode(node *ast_node.ASTNode) (*runtime_heap.VariableValue, error) {
	var functionReference = node.Body[0]

	if functionReference == nil {
		return nil, runtime_error.CreateError(
			"Cannot get reference to function",
			node,
		)
	}

	var functionVariable, funcionVariableErr = runtime.visitNode(functionReference)

	if funcionVariableErr != nil {
		return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
			"Cannot get function",
			functionReference,
		), funcionVariableErr)
	}

	var argumentsValues []*runtime_heap.VariableValue

	for _, argumentNode := range node.Arguments {
		var argumentValue, argumentValueErr = runtime.visitNode(argumentNode)

		if argumentValueErr != nil {
			return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
				"Cannot get function argument",
				argumentNode,
			), argumentValueErr)
		}

		argumentsValues = append(argumentsValues, argumentValue)
	}

	if functionVariable.ValueType == runtime_heap.TYPE_NATIVE_FUNCTION {
		runtime.callNativeFunction(functionVariable.NativeFunctionName, argumentsValues)
		return nil, nil
	}

	if functionVariable.ValueType != runtime_heap.TYPE_FUNCTION {
		return nil, runtime_error.CreateError(
			"Is not a function",
			functionReference,
		)
	}

	var innerRuntime = CreateRuntime(runtime.bridge)
	var argumentsNames = []string{}

	for _, funcParam := range functionVariable.FunctionValue.Params {
		if funcParam.Name == ast_node.AST_PARAM_FUNCTION_ARGUMENT_NAME {
			argumentsNames = append(argumentsNames, funcParam.Value)
		}
	}

	for index, argumentName := range argumentsNames {
		var createArgumentValueError = innerRuntime.heap.CreateVariable(argumentName)

		if createArgumentValueError != nil {
			return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
				"Cannot create variable for argument: "+argumentName,
				node,
			), createArgumentValueError)
		}

		if index < len(argumentsValues) {
			var argumentValue = argumentsValues[index]

			var setArgumentValueError = innerRuntime.heap.SetVariable(argumentName, argumentValue)

			if setArgumentValueError != nil {
				return nil, runtime_error.MergeRuntimeErrors(runtime_error.CreateError(
					"Cannot set value for argument: "+argumentName,
					node,
				), setArgumentValueError)
			}
		}
	}

	for _, functionBodyNode := range functionVariable.FunctionValue.Body {
		var bodyNodeValue, bodyNodeErr = innerRuntime.visitNode(functionBodyNode)

		if bodyNodeErr != nil {
			return nil, bodyNodeErr
		}

		if functionBodyNode.Code == ast_node.AST_NODE_CODE_RETURN {
			return bodyNodeValue, nil
		}
	}

	return nil, nil
}

func (runtime *Runtime) defineEnvByBridge() error {
	var printFuncName = "print"
	var definePrintVariableErr = runtime.heap.CreateVariable(printFuncName)

	if definePrintVariableErr != nil {
		return runtime_error.MergeRuntimeErrors(runtime_error.RuntimeError{
			Message: "Cannot create print env variable",
		}, definePrintVariableErr)
	}

	var defineNativePrintErr = runtime.heap.SetVariable(printFuncName, &runtime_heap.VariableValue{
		ValueType:          runtime_heap.TYPE_NATIVE_FUNCTION,
		NativeFunctionName: printFuncName,
	})

	if defineNativePrintErr != nil {
		return runtime_error.MergeRuntimeErrors(runtime_error.RuntimeError{
			Message: "Cannot define print env variable",
		}, defineNativePrintErr)
	}

	return nil
}

func (runtime *Runtime) callNativeFunction(name string, argumentsValues []*runtime_heap.VariableValue) error {
	if name == "print" {
		runtime.bridge.Print(argumentsValues...)
	}

	return nil
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
