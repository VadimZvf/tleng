package ast

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/ast_node_binary_expression"
	"github.com/VadimZvf/golang/ast_node_number"
	"github.com/VadimZvf/golang/ast_node_return"
	"github.com/VadimZvf/golang/ast_node_string"
	"github.com/VadimZvf/golang/ast_token_stream"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_keyword"
	"github.com/VadimZvf/golang/token_number"
	"github.com/VadimZvf/golang/token_read_property"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
	"github.com/VadimZvf/golang/token_variable_declaration"
)

type baseProcess = func(stream ast_node.ITokenStream, leftNode *ast_node.ASTNode) (resultNodes []*ast_node.ASTNode, err error)

type context struct {
	baseProcess baseProcess
}

func (ctx context) Process(stream ast_node.ITokenStream, currentCtx ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) (resultNodes []*ast_node.ASTNode, err error) {
	return ctx.baseProcess(stream, leftNode)
}

func createContext(baseProcess baseProcess) ast_node.IASTNodeProcessingContext {
	return context{
		baseProcess,
	}
}

func CreateAST(tokens []token.Token) (*ast_node.ASTNode, error) {
	var tokenStream = ast_token_stream.CreateTokenStream(tokens)
	var _, isEnd = tokenStream.Look()
	var ast = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
	}

	for !isEnd {
		var nodes, err = getNodes(&tokenStream, nil)

		for _, node := range nodes {
			ast_node.AppendNode(&ast, node)
		}

		if err != nil {
			return &ast, err
		}

		tokenStream.MoveNext()
		_, isEnd = tokenStream.Look()
	}

	return &ast, nil
}

func getNodes(stream ast_node.ITokenStream, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. At getNodes",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}
	var ctx = createContext(getNodes)

	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		return processVariableDeclaration(stream)

	case token_number.NUMBER:
		return ast_node_number.NumberProcessor(stream, ctx, nil)

	case token_string.STRING:
		return ast_node_string.StringProcessor(stream, ctx, nil)

	case token_keyword.KEY_WORD:
		return processKeyWord(stream)

	case token_function_declaration.FUNCTION_DECLARATION:
		return processFunction(stream)

	case token_return.RETURN_DECLARATION:
		return ast_node_return.ReturnProcessor(stream, ctx, nil)

	case token.OPEN_EXPRESSION:
		return processParenthesizedExpression(stream)

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		return ast_node_binary_expression.BinaryExpressionProcessor(stream, ctx, leftNode)

	case token.END_LINE:
		return []*ast_node.ASTNode{}, nil
	}

	return []*ast_node.ASTNode{}, parser_error.ParserError{
		Message:       "Unknown token. Code: " + currentToken.Code,
		StartPosition: currentToken.StartPosition,
		EndPosition:   currentToken.EndPosition,
	}
}

func processVariableDeclaration(stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at variable declaration processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var variableDeclarationNode = ast_node.CreateNode(currentToken)

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&variableDeclarationNode}, nil
	}

	switch nextToken.Code {
	case token.ASSIGNMENT:
		var variableNameParam = ast_node.GetVariableNameParam(&variableDeclarationNode)
		var referenceNode = ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_REFERENCE,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         variableNameParam.Value,
				StartPosition: variableNameParam.StartPosition,
				EndPosition:   variableNameParam.EndPosition,
			}},
			StartPosition: variableNameParam.StartPosition,
			EndPosition:   variableNameParam.EndPosition,
		}

		stream.MoveNext()

		var assignmentNodes, assignmentNodeParsingError = processAssignment(&referenceNode, stream)

		if assignmentNodeParsingError != nil {
			return []*ast_node.ASTNode{&variableDeclarationNode}, parser_error.MergeParserErrors(parser_error.ParserError{
				Message: "Parsing error. At assignment with variable declaration",
			}, assignmentNodeParsingError)
		}

		if len(assignmentNodes) != 1 {
			return []*ast_node.ASTNode{&variableDeclarationNode}, parser_error.ParserError{
				Message:       "Parsing error. Should assign only one node. But received: " + fmt.Sprint(len(assignmentNodes)),
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}
		}

		return []*ast_node.ASTNode{&variableDeclarationNode, assignmentNodes[0]}, nil
	}

	return []*ast_node.ASTNode{&variableDeclarationNode}, nil
}

func processKeyWord(stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at variable declaration processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var referenceNode = ast_node.CreateNode(currentToken)

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&referenceNode}, nil
	}

	switch nextToken.Code {
	case token.ASSIGNMENT:
		stream.MoveNext()
		return processAssignment(&referenceNode, stream)

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		return processBinaryExpression(&referenceNode, stream)

	case token_read_property.READ_PROPERTY:
		stream.MoveNext()
		return processReadProperty(&referenceNode, stream)

	case token.OPEN_EXPRESSION:
		stream.MoveNext()
		return processCallExpression(&referenceNode, stream)
	}

	return []*ast_node.ASTNode{&referenceNode}, nil
}

func processAssignment(leftNode *ast_node.ASTNode, stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at assignment processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var assignmentNode = ast_node.CreateNode(currentToken)

	stream.MoveNext()

	var rightNodes, rightNodesParsingError = getNodes(stream, nil)

	if rightNodesParsingError != nil {
		return []*ast_node.ASTNode{leftNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse right node of assignment expression",
		}, rightNodesParsingError)
	}

	if len(rightNodes) != 1 {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Parsing error. Right node of assignment expression should have only one node. But received: " + fmt.Sprint(len(rightNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	ast_node.AppendNode(&assignmentNode, leftNode)
	ast_node.AppendNodes(&assignmentNode, rightNodes)

	return []*ast_node.ASTNode{&assignmentNode}, nil
}

func processBinaryExpression(leftNode *ast_node.ASTNode, stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at binary expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var binaryNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()
	var rightNodes, rightNodeError = getNodes(stream, nil)

	if rightNodeError != nil {
		return []*ast_node.ASTNode{leftNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse right node of binary expression",
		}, rightNodeError)
	}

	if len(rightNodes) != 1 {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Parsing error. Right node of binary expression should have only one node. But received: " + fmt.Sprint(len(rightNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	ast_node.AppendNodes(&binaryNode, []*ast_node.ASTNode{
		leftNode,
		rightNodes[0],
	})

	return []*ast_node.ASTNode{&binaryNode}, nil
}

func processParenthesizedExpression(stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at parenthesized expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var parenthesizedExpressionNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()
	var valueNodes, valueNodeError = getNodes(stream, nil)

	if valueNodeError != nil {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse value node of parenthesized expression",
		}, valueNodeError)
	}

	if len(valueNodes) != 1 {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.ParserError{
			Message:       "Parsing error. Parenthesized expression should have only one value node. But received: " + fmt.Sprint(len(valueNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		ast_node.AppendNodes(&parenthesizedExpressionNode, valueNodes)

		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
	}

	if nextToken.Code != token.CLOSE_EXPRESSION {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.ParserError{
			Message:       "Failed parse parenthesized expression. Expression should be closed. But received token: " + nextToken.Code,
			StartPosition: nextToken.StartPosition,
			EndPosition:   nextToken.EndPosition,
		}
	}
	parenthesizedExpressionNode.EndPosition = nextToken.EndPosition

	stream.MoveNext()
	ast_node.AppendNodes(&parenthesizedExpressionNode, valueNodes)

	nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
	}

	switch nextToken.Code {
	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		var binaryExpressionNodes, binaryExpressionNodeError = processBinaryExpression(&parenthesizedExpressionNode, stream)

		if binaryExpressionNodeError != nil {
			return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.MergeParserErrors(parser_error.ParserError{
				Message: "Failed parse binary expression node in parenthesized expression",
			}, binaryExpressionNodeError)
		}

		if len(binaryExpressionNodes) != 1 {
			return binaryExpressionNodes, parser_error.ParserError{
				Message:       "Parsing error. Binary expression should has only one node. But received: " + fmt.Sprint(len(binaryExpressionNodes)),
				StartPosition: nextToken.StartPosition,
				EndPosition:   nextToken.EndPosition,
			}
		}

		return binaryExpressionNodes, nil

	case token.OPEN_EXPRESSION:
		stream.MoveNext()
		return processCallExpression(&parenthesizedExpressionNode, stream)
	}

	return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
}

func processFunction(stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at function processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var functionNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()

	var nextToken, isEndNext = stream.Look()

	if isEndNext {
		return []*ast_node.ASTNode{&functionNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Function should have body",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if nextToken.Code != token.OPEN_BLOCK {
		return []*ast_node.ASTNode{&functionNode}, parser_error.ParserError{
			Message:       "Function should have body",
			StartPosition: nextToken.StartPosition,
			EndPosition:   nextToken.EndPosition,
		}
	}

	stream.MoveNext()
	nextToken, isEndNext = stream.Look()

	for !isEndNext && nextToken.Code != token.CLOSE_BLOCK {
		var bodyNodes, bodyNodeParsingError = getNodes(stream, nil)

		if bodyNodeParsingError != nil {
			return []*ast_node.ASTNode{&functionNode}, parser_error.MergeParserErrors(parser_error.ParserError{
				Message: "Failed parsing in function body",
			}, bodyNodeParsingError)
		}

		ast_node.AppendNodes(&functionNode, bodyNodes)
		stream.MoveNext()
		nextToken, isEndNext = stream.Look()
	}

	functionNode.EndPosition = nextToken.EndPosition

	return []*ast_node.ASTNode{&functionNode}, nil
}

func processReturn(stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at return processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var returnNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()

	var valueNodes, valueNodeError = getNodes(stream, nil)

	if valueNodeError != nil {
		return []*ast_node.ASTNode{&returnNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse value node of return declaration",
		}, valueNodeError)
	}

	if len(valueNodes) != 1 {
		return []*ast_node.ASTNode{&returnNode}, parser_error.ParserError{
			Message:       "Parsing error. Return declaration should have only one value node. But received: " + fmt.Sprint(len(valueNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	ast_node.AppendNodes(&returnNode, valueNodes)

	return []*ast_node.ASTNode{&returnNode}, nil
}

func processReadProperty(leftNode *ast_node.ASTNode, stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at read property expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var nextReadPropetryNode = ast_node.CreateNode(currentToken)
	ast_node.AppendNode(&nextReadPropetryNode, leftNode)

	stream.MoveNext()
	var propertyToken, isEndAtProperty = stream.Look()

	if isEndAtProperty {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at read property expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	nextReadPropetryNode.Params = []ast_node.ASTNodeParam{
		{
			Name:          ast_node.AST_PARAM_PROPERTY_NAME,
			Value:         propertyToken.Value,
			StartPosition: propertyToken.StartPosition,
			EndPosition:   propertyToken.EndPosition,
		},
	}

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&nextReadPropetryNode}, nil
	}

	switch nextToken.Code {
	case token_read_property.READ_PROPERTY:
		stream.MoveNext()
		return processReadProperty(&nextReadPropetryNode, stream)

	case token.ASSIGNMENT:
		stream.MoveNext()
		return processAssignment(&nextReadPropetryNode, stream)

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		return processBinaryExpression(&nextReadPropetryNode, stream)

	case token.OPEN_EXPRESSION:
		stream.MoveNext()
		return processCallExpression(&nextReadPropetryNode, stream)
	}

	return []*ast_node.ASTNode{&nextReadPropetryNode}, nil
}

func processCallExpression(leftNode *ast_node.ASTNode, stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at call expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var callNode = ast_node.ASTNode{
		Code:          ast_node.AST_NODE_CODE_CALL_EXPRESSION,
		StartPosition: currentToken.StartPosition,
		EndPosition:   currentToken.EndPosition,
		Body:          []*ast_node.ASTNode{leftNode},
	}

	stream.MoveNext()

	var arguments, argumentsParsingError = processCallExpressionArguments(stream)

	callNode.Arguments = arguments

	if argumentsParsingError != nil {
		return []*ast_node.ASTNode{&callNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Something wrong at arguments call expression processing",
		}, argumentsParsingError)
	}

	var endCallToken, isEndAtEndCall = stream.Look()

	if isEndAtEndCall {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at call expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if endCallToken.Code != token.CLOSE_EXPRESSION {
		return []*ast_node.ASTNode{&callNode}, parser_error.ParserError{
			Message:       "Unknow token. Expected end call expression. But received: " + endCallToken.Code,
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	callNode.EndPosition = endCallToken.EndPosition

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&callNode}, nil
	}

	switch nextToken.Code {
	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		return processBinaryExpression(&callNode, stream)
	case token.OPEN_EXPRESSION:
		stream.MoveNext()
		return processCallExpression(&callNode, stream)
	}

	return []*ast_node.ASTNode{&callNode}, nil
}

func processCallExpressionArguments(stream ast_node.ITokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()
	var arguments = []*ast_node.ASTNode{}

	for !isEnd && currentToken.Code != token.CLOSE_EXPRESSION {
		var argument, argumentParsingError = getNodes(stream, nil)

		if argumentParsingError != nil {
			return arguments, parser_error.MergeParserErrors(parser_error.ParserError{
				Message: "Something wrong at call argument processing",
			}, argumentParsingError)
		}

		if len(argument) != 1 {
			return arguments, parser_error.ParserError{
				Message:       "Parsing error. Argument declaration should have only one value node. But received: " + fmt.Sprint(len(argument)),
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}
		}

		arguments = append(arguments, argument[0])

		stream.MoveNext()
		currentToken, isEnd = stream.Look()

		if currentToken.Code != token.COMMA && currentToken.Code != token.CLOSE_EXPRESSION {
			return arguments, parser_error.ParserError{
				Message:       "Parsing error. Argument declarations should devided by comma. But received: " + currentToken.Code,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}
		}

		if currentToken.Code == token.COMMA {
			stream.MoveNext()
			currentToken, isEnd = stream.Look()
		}
	}

	return arguments, nil
}
