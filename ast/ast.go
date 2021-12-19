package ast

import (
	"github.com/VadimZvf/golang/ast_factory"
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/ast_token_stream"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_keyword"
	"github.com/VadimZvf/golang/token_number"
	"github.com/VadimZvf/golang/token_variable_declaration"
)

func CreateAST(tokens []token.Token) *ast_node.ASTNode {
	var tokenStream = ast_token_stream.CreateTokenStream(tokens)
	var currentToken, isEnd = tokenStream.Look()
	var factory = ast_factory.CreateASTFactory()

	for !isEnd {
		switch currentToken.Code {
		case token_variable_declaration.VARIABLE_DECLARAION:
			var variableNode = createNode(currentToken)
			factory.Append(&variableNode)

			var nextToken, isEnd = tokenStream.LookNext()

			if isEnd {
				break
			}

			if nextToken.Code == token.ASSIGNMENT {
				var variableNameParam = token_variable_declaration.GetVariableNameParam(currentToken)
				var fakeKeyWordToken = token.Token{
					Code:          token_keyword.KEY_WORD,
					StartPosition: currentToken.StartPosition,
					EndPosition:   currentToken.EndPosition,
					Value:         variableNameParam.Value,
				}
				processKeyWordToken(fakeKeyWordToken, &tokenStream, &factory)
			}

		case token_keyword.KEY_WORD:
			processKeyWordToken(currentToken, &tokenStream, &factory)

		case token_function_declaration.FUNCTION_DECLARATION:
			processFunctionToken(currentToken, &tokenStream, &factory)

		case token.CLOSE_BLOCK:
			factory.MovePointerToParent()
		}

		tokenStream.MoveNext()
		currentToken, isEnd = tokenStream.Look()
	}

	return factory.GetAST()
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

func processKeyWordToken(currentToken token.Token, stream *ast_token_stream.TokenStream, factory *ast_factory.ASTFactory) {
	var nextToken, isEnd = stream.LookNext()

	if isEnd {
		return
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
			return
		}

		var valueNode = createNode(valueToken)

		assignmentNode.Body = []*ast_node.ASTNode{&valueNode}

		factory.Append(&assignmentNode)
		stream.MoveNext()
	}
}

func processFunctionToken(currentToken token.Token, stream *ast_token_stream.TokenStream, factory *ast_factory.ASTFactory) {
	var nextToken, isEnd = stream.LookNext()

	if isEnd {
		return
	}

	switch nextToken.Code {
	case token.OPEN_BLOCK:
		var functionNode = createNode(currentToken)
		factory.Append(&functionNode)
		factory.MovePointerLastNodeBody()
		stream.MoveNext()
	}
}
