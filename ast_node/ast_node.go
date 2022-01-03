package ast_node

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
