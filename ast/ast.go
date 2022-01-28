package ast

import (
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/ast_node_assignment"
	"github.com/VadimZvf/golang/ast_node_binary_expression"
	"github.com/VadimZvf/golang/ast_node_boolean"
	"github.com/VadimZvf/golang/ast_node_call_expression"
	"github.com/VadimZvf/golang/ast_node_function"
	"github.com/VadimZvf/golang/ast_node_number"
	"github.com/VadimZvf/golang/ast_node_parenthesized_expression"
	"github.com/VadimZvf/golang/ast_node_read_property"
	"github.com/VadimZvf/golang/ast_node_reference"
	"github.com/VadimZvf/golang/ast_node_return"
	"github.com/VadimZvf/golang/ast_node_string"
	"github.com/VadimZvf/golang/ast_node_variable_declaration"
	"github.com/VadimZvf/golang/ast_token_stream"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_boolean"
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
	var ctx = context{}

	for !isEnd {
		var nodes, err = ctx.Process(&tokenStream, ctx, nil)

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

type context struct{}

func (ctx context) Process(stream ast_node.ITokenStream, currentCtx ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) (resultNodes []*ast_node.ASTNode, err error) {
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
		return ast_node_variable_declaration.VariableDeclarationProcessor(stream, ctx, leftNode)

	case token_number.NUMBER:
		return ast_node_number.NumberProcessor(stream, ctx, leftNode)

	case token_string.STRING:
		return ast_node_string.StringProcessor(stream, ctx, leftNode)

	case token_boolean.BOOLEAN:
		return ast_node_boolean.BooleanProcessor(stream, ctx, leftNode)

	case token_keyword.KEY_WORD:
		return ast_node_reference.ReferenceProcessor(stream, ctx, leftNode)

	case token_function_declaration.FUNCTION_DECLARATION:
		return ast_node_function.FunctionProcessor(stream, ctx, leftNode)

	case token_return.RETURN_DECLARATION:
		return ast_node_return.ReturnProcessor(stream, ctx, leftNode)

	case token.OPEN_EXPRESSION:
		if leftNode == nil {
			return ast_node_parenthesized_expression.ParenthesizedExpressionProcessor(stream, ctx, leftNode)
		}

		return ast_node_call_expression.CallExpressionProcessor(stream, ctx, leftNode)

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		return ast_node_binary_expression.BinaryExpressionProcessor(stream, ctx, leftNode)

	case token_read_property.READ_PROPERTY:
		return ast_node_read_property.ReadPropertyProcessor(stream, ctx, leftNode)

	case token.ASSIGNMENT:
		return ast_node_assignment.AssignmentProcessor(stream, ctx, leftNode)

	case token.END_LINE:
		return []*ast_node.ASTNode{}, nil
	}

	return []*ast_node.ASTNode{}, parser_error.ParserError{
		Message:       "Unknown token. Code: " + currentToken.Code,
		StartPosition: currentToken.StartPosition,
		EndPosition:   currentToken.EndPosition,
	}
}
