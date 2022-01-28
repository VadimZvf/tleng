package ast_node_boolean

import (
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
)

var BooleanProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode != nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Left node not supported for boolean node",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at boolean processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var booleanNode = ast_node.CreateNode(currentToken)

	return []*ast_node.ASTNode{&booleanNode}, nil
}
