package ast

import (
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_keyword"
	"github.com/VadimZvf/golang/token_number"
	"github.com/VadimZvf/golang/token_variable_declaration"
)

type ASTNodeParam struct {
	Name  string
	Value string
	// Debug data
	StartPosition int
	EndPosition   int
}

type ASTNode struct {
	Code string
	// Detail information, like function arguments, return values etc...
	Params  []ASTNodeParam
	Child   *ASTNode
	Sibling *ASTNode
	// Debug data
	StartPosition int
	EndPosition   int
}

const AST_NODE_CODE_ROOT = "ROOT"
const AST_NODE_CODE_VARIABLE_DECLARATION = "VARIABLE_DECLARATION"
const AST_NODE_CODE_ASSIGNMENT = "ASSIGNMENT"
const AST_NODE_CODE_NUMBER = "NUMBER"
const AST_NODE_CODE_REFERENDE = "REFERENDE"
const AST_NODE_CODE_FUNCTION = "FUNCTION"

const AST_PARAM_VARIABLE_NAME = "VARIABLE_NAME"
const AST_PARAM_NUMBER_VALUE = "NUMBER_VALUE"

func CreateAST(tokens []token.Token) ASTNode {
	var root = ASTNode{
		Code: AST_NODE_CODE_ROOT,
	}

	var current *ASTNode = &root

	for _, currentToken := range tokens {
		switch currentToken.Code {
		case token_variable_declaration.VARIABLE_DECLARAION:
			var variableName = currentToken.Params[len(currentToken.Params)-1]
			var variableNode = ASTNode{
				Code: AST_NODE_CODE_VARIABLE_DECLARATION,
				Params: []ASTNodeParam{{
					Name:          AST_PARAM_VARIABLE_NAME,
					Value:         variableName.Value,
					StartPosition: variableName.StartPosition,
					EndPosition:   variableName.EndPosition,
				}},
				// Debug data
				StartPosition: current.StartPosition,
				EndPosition:   current.EndPosition,
			}

			current.Sibling = &variableNode
			current = &variableNode

		case token.ASSIGNMENT:
			if current.Code == AST_NODE_CODE_VARIABLE_DECLARATION || current.Code == AST_NODE_CODE_REFERENDE {
				var variableName = getVariableNameParam(*current)
				var assignmentNode = ASTNode{
					Code: AST_NODE_CODE_ASSIGNMENT,
					Params: []ASTNodeParam{{
						Name:          AST_PARAM_VARIABLE_NAME,
						Value:         variableName.Value,
						StartPosition: variableName.StartPosition,
						EndPosition:   variableName.EndPosition,
					}},
					// Debug data
					StartPosition: current.StartPosition,
					EndPosition:   current.EndPosition,
				}

				current.Sibling = &assignmentNode
				current = &assignmentNode
			}

		case token_number.NUMBER:
			var numberNode = ASTNode{
				Code: AST_NODE_CODE_NUMBER,
				Params: []ASTNodeParam{{
					Name:          AST_PARAM_NUMBER_VALUE,
					Value:         currentToken.Value,
					StartPosition: currentToken.StartPosition,
					EndPosition:   currentToken.EndPosition,
				}},
				// Debug data
				StartPosition: current.StartPosition,
				EndPosition:   current.EndPosition,
			}

			if current.Code == AST_NODE_CODE_ASSIGNMENT {
				current.Child = &numberNode
			}

		case token_keyword.KEY_WORD:
			var referenceNode = ASTNode{
				Code: AST_NODE_CODE_REFERENDE,
				Params: []ASTNodeParam{{
					Name:          AST_PARAM_VARIABLE_NAME,
					Value:         currentToken.Value,
					StartPosition: currentToken.StartPosition,
					EndPosition:   currentToken.EndPosition,
				}},
				// Debug data
				StartPosition: current.StartPosition,
				EndPosition:   current.EndPosition,
			}

			if current.Code == AST_NODE_CODE_ASSIGNMENT && current.Child == nil {
				current.Child = &referenceNode
			} else {
				current.Sibling = &referenceNode
				current = &referenceNode
			}
		}
	}

	return root
}

func getVariableNameParam(node ASTNode) ASTNodeParam {
	for _, param := range node.Params {
		if param.Name == AST_PARAM_VARIABLE_NAME {
			return param
		}
	}

	return ASTNodeParam{}
}
