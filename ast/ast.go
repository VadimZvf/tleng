package ast

import (
	"github.com/VadimZvf/golang/ast_error"
	"github.com/VadimZvf/golang/ast_factory"
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/ast_token_stream"
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
	var factory = ast_factory.CreateASTFactory()

	for !isEnd {
		var nodes, err = getNode(&tokenStream, &factory)

		if err != nil {
			return factory.GetAST(), err
		}

		tokenStream.MoveNext()
		_, isEnd = tokenStream.Look()

		if factory.GetLastInWorkStack() == nil {
			for _, node := range nodes {
				factory.Append(node)
			}
		}
	}

	return factory.GetAST(), nil
}

func getNode(stream *ast_token_stream.TokenStream, factory *ast_factory.ASTFactory) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, ast_error.AstError{
			Message: "File ended",
		}
	}

	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		var variableNode = ast_node.CreateNode(currentToken)

		var nextToken, isEndNext = stream.LookNext()

		if nextToken.Code != token.ASSIGNMENT || isEndNext {
			return []*ast_node.ASTNode{&variableNode}, nil
		}

		factory.PushToWorkStack(&variableNode)
		stream.MoveNext()
		var nextNodes, err = getNode(stream, factory)

		if err != nil {
			return []*ast_node.ASTNode{&variableNode}, err
		}

		return append([]*ast_node.ASTNode{&variableNode}, nextNodes...), nil

	case token.ASSIGNMENT:
		var referenceNode = factory.PopWorkStack()

		if referenceNode == nil {
			return []*ast_node.ASTNode{}, ast_error.AstError{
				Message: "Reference for assignment not defined! start: " + string(currentToken.StartPosition),
			}
		}

		var variableNameParam = ast_node.GetVariableNameParam(referenceNode)
		var assignmentNode = ast_node.CreateNode(currentToken)
		assignmentNode.Params = []ast_node.ASTNodeParam{{
			Name:          ast_node.AST_PARAM_VARIABLE_NAME,
			Value:         variableNameParam.Value,
			StartPosition: variableNameParam.StartPosition,
			EndPosition:   variableNameParam.EndPosition,
		}}

		stream.MoveNext()
		var valueNodes, err = getNode(stream, factory)

		if err != nil {
			return valueNodes, err
		}

		appendNodes(&assignmentNode, valueNodes)

		return []*ast_node.ASTNode{&assignmentNode}, nil

	case token_number.NUMBER:
		var numberNode = ast_node.CreateNode(currentToken)
		var nextToken, isEndNext = stream.LookNext()

		if nextToken.Code == token.END_LINE || isEndNext {
			stream.MoveNext()
			return []*ast_node.ASTNode{&numberNode}, nil
		}

		factory.PushToWorkStack(&numberNode)
		return []*ast_node.ASTNode{}, nil

	case token_string.STRING:
		var stringNode = ast_node.CreateNode(currentToken)
		var nextToken, isEndNext = stream.LookNext()

		if nextToken.Code == token.END_LINE || isEndNext {
			stream.MoveNext()
			return []*ast_node.ASTNode{&stringNode}, nil
		}

		factory.PushToWorkStack(&stringNode)
		return []*ast_node.ASTNode{}, nil

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		var leftNode = factory.PopWorkStack()
		stream.MoveNext()
		var binaryNode = ast_node.CreateNode(currentToken)
		stream.MoveNext()

		var rightNode, err = getNode(stream, factory)

		if err != nil {
			return []*ast_node.ASTNode{}, err
		}

		appendNode(&binaryNode, leftNode)
		appendNodes(&binaryNode, rightNode)

		return []*ast_node.ASTNode{&binaryNode}, nil

	case token_read_property.READ_PROPERTY:
		var leftNode = factory.PopWorkStack()
		stream.MoveNext()
		var readNode = ast_node.CreateNode(currentToken)
		stream.MoveNext()

		var rightNodes, err = getNode(stream, factory)

		if err != nil {
			return []*ast_node.ASTNode{}, err
		}

		appendNode(&readNode, leftNode)
		appendNodes(&readNode, rightNodes)

		return []*ast_node.ASTNode{&readNode}, nil

	case token_return.RETURN_DECLARATION:
		var returnNode = ast_node.CreateNode(currentToken)
		stream.MoveNext()
		var returnValueNode, err = getNode(stream, factory)

		if err != nil {
			return []*ast_node.ASTNode{&returnNode}, err
		}

		appendNodes(&returnNode, returnValueNode)
		return []*ast_node.ASTNode{&returnNode}, nil

	case token_keyword.KEY_WORD:
		var keyWordNode = processKeyWordToken(stream, factory)
		return []*ast_node.ASTNode{&keyWordNode}, nil

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionNode = processFunctionToken(stream, factory)
		return []*ast_node.ASTNode{&functionNode}, nil

	case token.OPEN_EXPRESSION:
		var parenthesizedExpressionNode = ast_node.CreateNode(currentToken)
		factory.PushToWorkStack(&parenthesizedExpressionNode)

		return []*ast_node.ASTNode{}, nil

	case token.CLOSE_EXPRESSION:
		var childNode = factory.PopWorkStack()

		if childNode == nil {
			return []*ast_node.ASTNode{}, nil
		}

		if childNode.Code == ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION {
			return []*ast_node.ASTNode{childNode}, nil
		}

		var parenthesizedExpressionNode = factory.PopWorkStack()
		parenthesizedExpressionNode.Body = []*ast_node.ASTNode{childNode}

		return []*ast_node.ASTNode{parenthesizedExpressionNode}, nil

	case token.END_LINE:
		var lastNode = factory.PopWorkStack()
		return []*ast_node.ASTNode{lastNode}, nil
	}

	return []*ast_node.ASTNode{}, nil
}

func processKeyWordToken(stream *ast_token_stream.TokenStream, factory *ast_factory.ASTFactory) ast_node.ASTNode {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}
	}

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return ast_node.CreateNode(currentToken)
	}

	switch nextToken.Code {
	case token.ASSIGNMENT:
		var assignmentNode = ast_node.CreateNode(nextToken)
		assignmentNode.Params = []ast_node.ASTNodeParam{{
			Name:          ast_node.AST_PARAM_VARIABLE_NAME,
			Value:         currentToken.Value,
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}}

		stream.MoveNext()
		stream.MoveNext()

		var valueNodes, error = getNode(stream, factory)

		if error != nil {
			return ast_node.ASTNode{}
		}

		appendNodes(&assignmentNode, valueNodes)

		return assignmentNode

	case token.OPEN_EXPRESSION:
		stream.MoveNext()
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_FUNCTION_CALL,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_FUNCTION_NAME,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	return ast_node.CreateNode(currentToken)
}

func processFunctionToken(stream *ast_token_stream.TokenStream, factory *ast_factory.ASTFactory) ast_node.ASTNode {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}
	}

	var functionNode = ast_node.CreateNode(currentToken)

	stream.MoveNext()
	currentToken, isEnd = stream.Look()

	for !isEnd && currentToken.Code != token.CLOSE_BLOCK {
		stream.MoveNext()
		currentToken, isEnd = stream.Look()

		var bodyNodes, err = getNode(stream, factory)

		if err != nil {
			return functionNode
		}

		appendNodes(&functionNode, bodyNodes)
	}

	return functionNode
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
