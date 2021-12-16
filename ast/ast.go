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
	var currentToken, isEnd = tokenStream.Next()
	var factory = ast_factory.CreateASTFactory()

	var stack []*ast_node.ASTNode = *&[]*ast_node.ASTNode{}

	for !isEnd {
		switch currentToken.Code {
		case token_variable_declaration.VARIABLE_DECLARAION:
			var variableNode = createNode(currentToken)
			factory.Append(&variableNode)

		case token.ASSIGNMENT:
			if isKeyWordLastInStack(stack) {
				var lastStackNode = stack[0]
				stack = stack[1:]

				if lastStackNode.Code == ast_node.AST_NODE_CODE_REFERENDE {
					var variableName = getVariableNameParam(lastStackNode)
					var assignmentNode = ast_node.ASTNode{
						Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
						Params: []ast_node.ASTNodeParam{{
							Name:          ast_node.AST_PARAM_VARIABLE_NAME,
							Value:         variableName.Value,
							StartPosition: variableName.StartPosition,
							EndPosition:   variableName.EndPosition,
						}},
						// Debug data
						StartPosition: lastStackNode.StartPosition,
						EndPosition:   lastStackNode.EndPosition,
					}

					factory.Append(&assignmentNode)
				}
				break
			}

			if factory.GetCurrent().Code == ast_node.AST_NODE_CODE_VARIABLE_DECLARATION || factory.GetCurrent().Code == ast_node.AST_NODE_CODE_REFERENDE {
				var variableName = getVariableNameParam(factory.GetCurrent())
				var assignmentNode = ast_node.ASTNode{
					Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
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

				factory.Append(&assignmentNode)
			}

		case token_number.NUMBER:
			var numberNode = createNode(currentToken)

			if factory.GetCurrent().Code == ast_node.AST_NODE_CODE_ASSIGNMENT {
				factory.Append(&numberNode)
			}

		case token_keyword.KEY_WORD:
			var referenceNode = createNode(currentToken)

			if factory.GetCurrent().Code == ast_node.AST_NODE_CODE_ASSIGNMENT {
				factory.Append(&referenceNode)
			} else {
				stack = append(stack, &referenceNode)
			}

		case token_function_declaration.FUNCTION_DECLARATION:
			var functionNode = createNode(currentToken)

			factory.Append(&functionNode)
			stack = append(stack, &functionNode)

		case token.OPEN_BLOCK:
			var functionName = token_function_declaration.GetFunctionNameParam(currentToken)

			var referenceNode = ast_node.ASTNode{
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

			factory.Append(&referenceNode)
			stack = append(stack, &referenceNode)
		}

		currentToken, isEnd = tokenStream.Next()
	}

	return factory.GetAST()
}

func isKeyWordLastInStack(stack []*ast_node.ASTNode) bool {
	if len(stack) == 0 {
		return false
	}

	var lastStackNode = stack[0]

	return lastStackNode.Code == ast_node.AST_NODE_CODE_REFERENDE
}

func getVariableNameParam(node *ast_node.ASTNode) ast_node.ASTNodeParam {
	for _, param := range node.Params {
		if param.Name == ast_node.AST_PARAM_VARIABLE_NAME {
			return param
		}
	}

	return ast_node.ASTNodeParam{}
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