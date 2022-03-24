package ast_node_function

import (
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var FunctionProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode != nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Left node not supported for function node",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at function processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var functionNode = ast_node.CreateNode(currentToken)
	stream.MoveNext()

	var nextToken, isEndNext = stream.Look()

	if isEndNext {
		return []*ast_node.ASTNode{&functionNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Function should have body",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if nextToken.Code != token.OPEN_BLOCK {
		return []*ast_node.ASTNode{&functionNode}, parser_error.ParserError{
			Message:       "Function should have body",
			StartPosition: nextToken.StartPosition,
			EndPosition:   nextToken.EndPosition,
		}
	}

	var bodyNodes, bodyNodeParsingError = context.Process(stream, context, nil)

	if bodyNodeParsingError != nil {
		return []*ast_node.ASTNode{&functionNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Failed parsing in function body",
		}, bodyNodeParsingError)
	}

	ast_node.AppendNodes(&functionNode, bodyNodes)
	nextToken, _ = stream.Look()

	functionNode.EndPosition = nextToken.EndPosition

	return []*ast_node.ASTNode{&functionNode}, nil
}
