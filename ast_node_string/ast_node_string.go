package ast_node_string

import (
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
)

var StringProcessor ast_node.ASTNodeProcessor = proccess

func proccess(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode != nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Left node not supported for string node",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at number processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var stringNode = ast_node.CreateNode(currentToken)

	if !ast_node.IsNextArithmeticToken(stream) {
		return []*ast_node.ASTNode{&stringNode}, nil
	}

	stream.MoveNext()

	return context.Process(stream, context, &stringNode)
}
