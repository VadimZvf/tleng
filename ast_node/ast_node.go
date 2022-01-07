package ast_node

import (
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_keyword"
	"github.com/VadimZvf/golang/token_number"
	"github.com/VadimZvf/golang/token_read_property"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
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
	Params    []ASTNodeParam
	Arguments []*ASTNode
	Body      []*ASTNode
	// Debug data
	StartPosition int
	EndPosition   int
}

const AST_NODE_CODE_ROOT = "ROOT"
const AST_NODE_CODE_VARIABLE_DECLARATION = "VARIABLE_DECLARATION"
const AST_NODE_CODE_ASSIGNMENT = "ASSIGNMENT"
const AST_NODE_CODE_BINARY_EXPRESSION = "BINARY_EXPRESSION"
const AST_NODE_CODE_PARENTHESIZED_EXPRESSION = "PARENTHESIZED_EXPRESSION"
const AST_NODE_CODE_CALL_EXPRESSION = "CALL_EXPRESSION"
const AST_NODE_CODE_READ_PROP = "READ_PROP"
const AST_NODE_CODE_NUMBER = "NUMBER"
const AST_NODE_CODE_STRING = "STRING"
const AST_NODE_CODE_REFERENCE = "REFERENCE"
const AST_NODE_CODE_FUNCTION = "FUNCTION"
const AST_NODE_CODE_RETURN = "RETURN"
const AST_NODE_CODE_FUNCTION_CALL = "FUNCTION_CALL"

const AST_PARAM_VARIABLE_NAME = "VARIABLE_NAME"
const AST_PARAM_FUNCTION_NAME = "FUNCTION_NAME"
const AST_PARAM_NUMBER_VALUE = "NUMBER_VALUE"
const AST_PARAM_STRING_VALUE = "STRING_VALUE"
const AST_PARAM_BINARY_EXPRESSION_TYPE = "BINARY_EXPRESSION_TYPE"
const AST_PARAM_PROPERTY_NAME = "PROPERTY_NAME"

func GetVariableNameParam(node *ASTNode) *ASTNodeParam {
	return GetParam(node, AST_PARAM_VARIABLE_NAME)
}

func GetNumberValueParam(node *ASTNode) *ASTNodeParam {
	return GetParam(node, AST_PARAM_NUMBER_VALUE)
}

func GetFunctionNameParam(node *ASTNode) *ASTNodeParam {
	return GetParam(node, AST_PARAM_FUNCTION_NAME)
}

func GetStringValueParam(node *ASTNode) *ASTNodeParam {
	return GetParam(node, AST_PARAM_STRING_VALUE)
}

func GetBinaryExpressionTypeParam(node *ASTNode) *ASTNodeParam {
	return GetParam(node, AST_PARAM_BINARY_EXPRESSION_TYPE)
}

func GetParam(node *ASTNode, paramCode string) *ASTNodeParam {
	for _, param := range node.Params {
		if param.Name == paramCode {
			return &param
		}
	}

	return nil
}

func CreateNode(currentToken token.Token) ASTNode {
	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		var variableName = token_variable_declaration.GetVariableNameParam(currentToken)
		return ASTNode{
			Code: AST_NODE_CODE_VARIABLE_DECLARATION,
			Params: []ASTNodeParam{{
				Name:          AST_PARAM_VARIABLE_NAME,
				Value:         variableName.Value,
				StartPosition: variableName.StartPosition,
				EndPosition:   variableName.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.ASSIGNMENT:
		return ASTNode{
			Code: AST_NODE_CODE_ASSIGNMENT,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		return ASTNode{
			Code: AST_NODE_CODE_BINARY_EXPRESSION,
			Params: []ASTNodeParam{{
				Name:          AST_PARAM_BINARY_EXPRESSION_TYPE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_read_property.READ_PROPERTY:
		return ASTNode{
			Code: AST_NODE_CODE_READ_PROP,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_number.NUMBER:
		return ASTNode{
			Code: AST_NODE_CODE_NUMBER,
			Params: []ASTNodeParam{{
				Name:          AST_PARAM_NUMBER_VALUE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_string.STRING:
		return ASTNode{
			Code: AST_NODE_CODE_STRING,
			Params: []ASTNodeParam{{
				Name:          AST_PARAM_STRING_VALUE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_keyword.KEY_WORD:
		return ASTNode{
			Code: AST_NODE_CODE_REFERENCE,
			Params: []ASTNodeParam{{
				Name:          AST_PARAM_VARIABLE_NAME,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_return.RETURN_DECLARATION:
		return ASTNode{
			Code: AST_NODE_CODE_RETURN,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionName = token_function_declaration.GetFunctionNameParam(currentToken)

		return ASTNode{
			Code: AST_NODE_CODE_FUNCTION,
			Params: []ASTNodeParam{{
				Name:          AST_PARAM_FUNCTION_NAME,
				Value:         functionName.Value,
				StartPosition: functionName.StartPosition,
				EndPosition:   functionName.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.OPEN_EXPRESSION:
		return ASTNode{
			Code: AST_NODE_CODE_PARENTHESIZED_EXPRESSION,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	default:
		return ASTNode{}
	}
}
