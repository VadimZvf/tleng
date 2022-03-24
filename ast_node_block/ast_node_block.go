package ast_node_block

import (
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var BlockProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode != nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected left node for block node",
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

	var blockNode = ast_node.CreateNode(currentToken)

	stream.MoveNext()

	var nextToken, isEndNext = stream.Look()

	if isEndNext {
		return []*ast_node.ASTNode{&blockNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Block should be closed",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	for !isEndNext && nextToken.Code != token.CLOSE_BLOCK {
		var bodyNodes, bodyNodeParsingError = context.Process(stream, context, nil)

		if bodyNodeParsingError != nil {
			return []*ast_node.ASTNode{&blockNode}, parser_error.MergeParserErrors(parser_error.ParserError{
				Message: "Failed parsing in function body",
			}, bodyNodeParsingError)
		}

		ast_node.AppendNodes(&blockNode, bodyNodes)
		stream.MoveNext()
		nextToken, isEndNext = stream.Look()
	}

	blockNode.EndPosition = nextToken.EndPosition

	return []*ast_node.ASTNode{&blockNode}, nil
}
