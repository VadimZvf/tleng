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
		var node, err = getNode(&tokenStream, &factory)
		factory.Append(&node)

		if err != nil {
			return factory.GetAST(), err
		}

		tokenStream.MoveNext()
		_, isEnd = tokenStream.Look()
	}

	return factory.GetAST(), nil
}

func getNode(stream *ast_token_stream.TokenStream, factory *ast_factory.ASTFactory) (ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}, ast_error.AstError{
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
			factory.Append(&variableNode)
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
			var valueNode, err = getNode(stream, factory)

			if err != nil {
				return valueNode, err
			}

			assignmentNode.Body = []*ast_node.ASTNode{&valueNode}

			stream.MoveNext()
			return assignmentNode, nil
		} else {
			return variableNode, nil
		}

	case token_number.NUMBER:
		return createNode(currentToken), nil

	case token.ADD:
		return createNode(currentToken), nil

	case token.SUBTRACT:
		return createNode(currentToken), nil

	case token_keyword.KEY_WORD:
		return processKeyWordToken(currentToken, stream), nil

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionNode = processFunctionToken(currentToken, stream)
		factory.MovePointerLastNodeBody()
		return functionNode, nil

	case token.CLOSE_BLOCK:
		factory.MovePointerToParent()

	case token.END_LINE:
		// Skip node
		// stream.MoveNext()

	}

	return ast_node.ASTNode{}, nil
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

func processKeyWordToken(currentToken token.Token, stream *ast_token_stream.TokenStream) ast_node.ASTNode {
	var nextToken, isEnd = stream.LookNext()

	if isEnd {
		return ast_node.ASTNode{}
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
		var valueToken, isEnd = stream.LookNext()

		if isEnd {
			return ast_node.ASTNode{}
		}

		var valueNode = createNode(valueToken)

		assignmentNode.Body = []*ast_node.ASTNode{&valueNode}

		stream.MoveNext()

		return assignmentNode
	}

	return ast_node.ASTNode{}
}

func processFunctionToken(currentToken token.Token, stream *ast_token_stream.TokenStream) ast_node.ASTNode {
	var nextToken, isEnd = stream.LookNext()

	if isEnd {
		return ast_node.ASTNode{}
	}

	switch nextToken.Code {
	case token.OPEN_BLOCK:
		var functionNode = createNode(currentToken)
		stream.MoveNext()

		return functionNode
	}

	return ast_node.ASTNode{}
}
