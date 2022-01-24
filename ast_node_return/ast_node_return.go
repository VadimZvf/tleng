package ast_node_return

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
)

var ReturnProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode != nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Left node not supported for return node",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at return processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var returnNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()

	var valueNodes, valueNodeError = context.Process(stream, context, nil)

	if valueNodeError != nil {
		return []*ast_node.ASTNode{&returnNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse value node of return declaration",
		}, valueNodeError)
	}

	if len(valueNodes) != 1 {
		return []*ast_node.ASTNode{&returnNode}, parser_error.ParserError{
			Message:       "Parsing error. Return declaration should have only one value node. But received: " + fmt.Sprint(len(valueNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	ast_node.AppendNodes(&returnNode, valueNodes)

	return []*ast_node.ASTNode{&returnNode}, nil
}
