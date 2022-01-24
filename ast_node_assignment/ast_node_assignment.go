package ast_node_assignment

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
)

var AssignmentProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode == nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Expected left node for assignment node",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at assignment processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var assignmentNode = ast_node.CreateNode(currentToken)

	stream.MoveNext()

	var rightNodes, rightNodesParsingError = context.Process(stream, context, nil)

	if rightNodesParsingError != nil {
		return []*ast_node.ASTNode{leftNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parse right node of assignment expression",
		}, rightNodesParsingError)
	}

	if len(rightNodes) != 1 {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Parsing error. Right node of assignment expression should have only one node. But received: " + fmt.Sprint(len(rightNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	ast_node.AppendNode(&assignmentNode, leftNode)
	ast_node.AppendNodes(&assignmentNode, rightNodes)

	return []*ast_node.ASTNode{&assignmentNode}, nil
}
