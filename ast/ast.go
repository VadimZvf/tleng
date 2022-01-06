package ast

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
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

func CreateAST(tokens []token.Token) (*ast_node.ASTNode, error) {
	var tokenStream = ast_token_stream.CreateTokenStream(tokens)
	var _, isEnd = tokenStream.Look()
	var ast = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
	}

	for !isEnd {
		var nodes, err = getNodes(&tokenStream)

		for _, node := range nodes {
			appendNode(&ast, node)
		}

		if err != nil {
			return &ast, err
		}

		tokenStream.MoveNext()
		_, isEnd = tokenStream.Look()
	}

	return &ast, nil
}

func getNodes(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. At getNodes",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		return processVariableDeclaration(stream)

	case token_number.NUMBER:
		return processNumber(stream)

	case token_string.STRING:
		return processString(stream)

	case token_keyword.KEY_WORD:
		return processKeyWord(stream)

	case token_function_declaration.FUNCTION_DECLARATION:
		return processFunction(stream)

	case token_return.RETURN_DECLARATION:
		return processReturn(stream)

	case token.OPEN_EXPRESSION:
		return processParenthesizedExpression(stream)

	case token.END_LINE:
		return []*ast_node.ASTNode{}, nil
	}

	return []*ast_node.ASTNode{}, parser_error.ParserError{
		Message:       "Unknown token. Code: " + currentToken.Code,
		StartPosition: currentToken.StartPosition,
		EndPosition:   currentToken.EndPosition,
	}
}

func processVariableDeclaration(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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
			return []*ast_node.ASTNode{&variableDeclarationNode}, mergeParserErrors(parser_error.ParserError{
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

func processKeyWord(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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

func processNumber(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at number processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var numberNode = ast_node.CreateNode(currentToken)

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&numberNode}, nil
	}

	switch nextToken.Code {
	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		return processBinaryExpression(&numberNode, stream)
	}

	return []*ast_node.ASTNode{&numberNode}, nil
}

func processString(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at number processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var stringNode = ast_node.CreateNode(currentToken)

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&stringNode}, nil
	}

	switch nextToken.Code {
	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		return processBinaryExpression(&stringNode, stream)
	}

	return []*ast_node.ASTNode{&stringNode}, nil
}

func processAssignment(leftNode *ast_node.ASTNode, stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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

	var rightNodes, rightNodesParsingError = getNodes(stream)

	if rightNodesParsingError != nil {
		return []*ast_node.ASTNode{leftNode}, mergeParserErrors(parser_error.ParserError{
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

	appendNode(&assignmentNode, leftNode)
	appendNodes(&assignmentNode, rightNodes)

	return []*ast_node.ASTNode{&assignmentNode}, nil
}

func processBinaryExpression(leftNode *ast_node.ASTNode, stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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
	var rightNodes, rightNodeError = getNodes(stream)

	if rightNodeError != nil {
		return []*ast_node.ASTNode{leftNode}, mergeParserErrors(parser_error.ParserError{
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

	appendNodes(&binaryNode, []*ast_node.ASTNode{
		leftNode,
		rightNodes[0],
	})

	return []*ast_node.ASTNode{&binaryNode}, nil
}

func processParenthesizedExpression(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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
	var valueNodes, valueNodeError = getNodes(stream)

	if valueNodeError != nil {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, mergeParserErrors(parser_error.ParserError{
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
		appendNodes(&parenthesizedExpressionNode, valueNodes)

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
	appendNodes(&parenthesizedExpressionNode, valueNodes)

	nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
	}

	switch nextToken.Code {
	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		stream.MoveNext()
		var binaryExpressionNodes, binaryExpressionNodeError = processBinaryExpression(&parenthesizedExpressionNode, stream)

		if binaryExpressionNodeError != nil {
			return []*ast_node.ASTNode{&parenthesizedExpressionNode}, mergeParserErrors(parser_error.ParserError{
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
	}

	return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
}

func processFunction(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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
		var bodyNodes, bodyNodeParsingError = getNodes(stream)

		if bodyNodeParsingError != nil {
			return []*ast_node.ASTNode{&functionNode}, mergeParserErrors(parser_error.ParserError{
				Message: "Failed parsing in function body",
			}, bodyNodeParsingError)
		}

		appendNodes(&functionNode, bodyNodes)
		stream.MoveNext()
		nextToken, isEndNext = stream.Look()
	}

	functionNode.EndPosition = nextToken.EndPosition

	return []*ast_node.ASTNode{&functionNode}, nil
}

func processReturn(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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

	var valueNodes, valueNodeError = getNodes(stream)

	if valueNodeError != nil {
		return []*ast_node.ASTNode{&returnNode}, mergeParserErrors(parser_error.ParserError{
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

	appendNodes(&returnNode, valueNodes)

	return []*ast_node.ASTNode{&returnNode}, nil
}

func processReadProperty(leftNode *ast_node.ASTNode, stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at read property expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var nextReadPropetryNode = ast_node.CreateNode(currentToken)
	appendNode(&nextReadPropetryNode, leftNode)

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

func processCallExpression(leftNode *ast_node.ASTNode, stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
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
		return []*ast_node.ASTNode{&callNode}, mergeParserErrors(parser_error.ParserError{
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

	return []*ast_node.ASTNode{&callNode}, nil
}

func processCallExpressionArguments(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()
	var arguments = []*ast_node.ASTNode{}

	for !isEnd && currentToken.Code != token.CLOSE_EXPRESSION {
		var argument, argumentParsingError = getNodes(stream)

		if argumentParsingError != nil {
			return arguments, mergeParserErrors(parser_error.ParserError{
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

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
// ┃         Utilities          ┃
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

func appendNodes(node *ast_node.ASTNode, children []*ast_node.ASTNode) {
	for _, child := range children {
		node.Body = append(node.Body, child)
	}
}

func appendNode(node *ast_node.ASTNode, child *ast_node.ASTNode) {
	node.Body = append(node.Body, child)
}

func mergeParserErrors(first error, second error) error {
	firstParserError, firstCastOk := first.(parser_error.ParserError)

	if !firstCastOk {
		return firstParserError
	}

	secondParserError, secondCastOk := second.(parser_error.ParserError)

	if !secondCastOk {
		return firstParserError
	}

	firstParserError.Message = secondParserError.Message + "\n  " + firstParserError.Message
	firstParserError.StartPosition = secondParserError.StartPosition
	firstParserError.EndPosition = secondParserError.EndPosition

	return firstParserError
}
