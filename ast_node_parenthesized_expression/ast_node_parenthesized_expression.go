package ast_node_parenthesized_expression

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var ParenthesizedExpressionProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at parenthesized expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var parenthesizedExpressionNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()
	var valueNodes, valueNodeError = context.Process(stream, context, nil)

	if valueNodeError != nil {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse value node of parenthesized expression",
		}, valueNodeError)
	}

	if len(valueNodes) != 1 {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.ParserError{
			Message:       "Parsing error. Parenthesized expression should have only one value node. But received: " + fmt.Sprint(len(valueNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		ast_node.AppendNodes(&parenthesizedExpressionNode, valueNodes)

		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
	}

	if nextToken.Code != token.CLOSE_EXPRESSION {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, parser_error.ParserError{
			Message:       "Failed parse parenthesized expression. Expression should be closed. But received token: " + nextToken.Code,
			StartPosition: nextToken.StartPosition,
			EndPosition:   nextToken.EndPosition,
		}
	}
	parenthesizedExpressionNode.EndPosition = nextToken.EndPosition
	ast_node.AppendNodes(&parenthesizedExpressionNode, valueNodes)

	stream.MoveNext()

	if !ast_node.IsNextExpressionToken(stream) {
		return []*ast_node.ASTNode{&parenthesizedExpressionNode}, nil
	}
	stream.MoveNext()

	return context.Process(stream, context, &parenthesizedExpressionNode)
}
