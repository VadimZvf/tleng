package ast_node_binary_expression

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
)

var BinaryExpressionProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode == nil {
		return []*ast_node.ASTNode{leftNode}, parser_error.CreateError(
			"Binary expression expect left node",
			currentToken.StartPosition,
			currentToken.EndPosition,
		)
	}

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at binary expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var binaryNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()
	var rightNodes, rightNodeError = context.Process(stream, context, nil)

	if rightNodeError != nil {
		return []*ast_node.ASTNode{leftNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse right node of binary expression",
		}, rightNodeError)
	}

	if len(rightNodes) != 1 {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Parsing error. Right node of binary expression should have only one node. But received: " + fmt.Sprint(len(rightNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	ast_node.AppendNodes(&binaryNode, []*ast_node.ASTNode{
		leftNode,
		rightNodes[0],
	})

	return []*ast_node.ASTNode{&binaryNode}, nil
}
