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
	"github.com/VadimZvf/golang/token_variable_declaration"
)

func CreateAST(tokens []token.Token) (*ast_node.ASTNode, error) {
	var tokenStream = ast_token_stream.CreateTokenStream(tokens)
	var _, isEnd = tokenStream.Look()
	var factory = ast_factory.CreateASTFactory()

	for !isEnd {
		var nodes, err = getNodes(&tokenStream)

		for _, child := range nodes {
			factory.Append(child)
		}

		if err != nil {
			return factory.GetAST(), err
		}

		tokenStream.MoveNext()
		_, isEnd = tokenStream.Look()
	}

	return factory.GetAST(), nil
}

func getNodes(stream *ast_token_stream.TokenStream) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, ast_error.AstError{
			Message: "File ended",
		}
	}

	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		var variableNode = createNode(currentToken)

		var nextToken, isEnd = stream.LookNext()

		if isEnd {
			break
		}

		if nextToken.Code == token.ASSIGNMENT {
			stream.MoveNext()

			var variableNameParam = token_variable_declaration.GetVariableNameParam(currentToken)
			var assignmentNode = createNode(nextToken)
			assignmentNode.Params = []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         variableNameParam.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}}

			stream.MoveNext()
			var valueNodes, err = getNodes(stream)

			if err != nil {
				return valueNodes, err
			}

			assignmentNode.Body = valueNodes

			stream.MoveNext()
			return []*ast_node.ASTNode{&variableNode, &assignmentNode}, nil
		} else {
			return []*ast_node.ASTNode{&variableNode}, nil
		}

	case token_number.NUMBER:
		var numberNode = createNode(currentToken)
		return []*ast_node.ASTNode{&numberNode}, nil

	case token.ADD:
		var addNode = createNode(currentToken)
		return []*ast_node.ASTNode{&addNode}, nil

	case token.SUBTRACT:
		var subtractNode = createNode(currentToken)
		return []*ast_node.ASTNode{&subtractNode}, nil

	case token_keyword.KEY_WORD:
		var keyWordNode = processKeyWordToken(stream)
		return []*ast_node.ASTNode{&keyWordNode}, nil

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionNode = processFunctionToken(stream)
		return []*ast_node.ASTNode{&functionNode}, nil

	case token.END_LINE:
		// Skip node
		// stream.MoveNext()

	}

	return []*ast_node.ASTNode{}, nil
}

func processKeyWordToken(stream *ast_token_stream.TokenStream) ast_node.ASTNode {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}
	}

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return createNode(currentToken)
	}

	switch nextToken.Code { 
	case token.ASSIGNMENT:
		var assignmentNode = createNode(nextToken)
		assignmentNode.Params = []ast_node.ASTNodeParam{{
			Name:          ast_node.AST_PARAM_VARIABLE_NAME,
			Value:         currentToken.Value,
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}}

		stream.MoveNext()
		stream.MoveNext()

		var valueNodes, error = getNodes(stream);

		if error != nil {
			return ast_node.ASTNode{}
		}

		assignmentNode.Body = valueNodes

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

	return createNode(currentToken)
}

func processFunctionToken(stream *ast_token_stream.TokenStream) ast_node.ASTNode {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}
	}

	var functionNode = createNode(currentToken)

	stream.MoveNext()
	currentToken, isEnd = stream.Look()

	for !isEnd && currentToken.Code != token.CLOSE_BLOCK {
		stream.MoveNext()
		currentToken, isEnd = stream.Look()

		var nodes, err = getNodes(stream)

		if err != nil {
			return functionNode
		}

		appendNodes(&functionNode, nodes)
	}

	return functionNode
}

func appendNodes(node *ast_node.ASTNode, children []*ast_node.ASTNode)  {
	for _, child := range children {
		node.Body = append(node.Body, child)
	}
}

func createNode(currentToken token.Token) ast_node.ASTNode {
	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		var variableName = token_variable_declaration.GetVariableNameParam(currentToken)
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         variableName.Value,
				StartPosition: variableName.StartPosition,
				EndPosition:   variableName.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.ASSIGNMENT:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.ADD:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_ADD,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.SUBTRACT:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_SUBTRACT,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_number.NUMBER:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_NUMBER,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_NUMBER_VALUE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_keyword.KEY_WORD:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_REFERENDE,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionName = token_function_declaration.GetFunctionNameParam(currentToken)

		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_FUNCTION,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_FUNCTION_NAME,
				Value:         functionName.Value,
				StartPosition: functionName.StartPosition,
				EndPosition:   functionName.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	default:
		return ast_node.ASTNode{}
	}
}
