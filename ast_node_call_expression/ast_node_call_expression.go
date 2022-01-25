package ast_node_call_expression

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var CallExpressionProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at call expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var callNode = ast_node.ASTNode{
		Code:          ast_node.AST_NODE_CODE_CALL_EXPRESSION,
		StartPosition: currentToken.StartPosition,
		EndPosition:   currentToken.EndPosition,
		Body:          []*ast_node.ASTNode{leftNode},
	}

	stream.MoveNext()

	var arguments, argumentsParsingError = processCallExpressionArguments(stream, context)

	callNode.Arguments = arguments

	if argumentsParsingError != nil {
		return []*ast_node.ASTNode{&callNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Something wrong at arguments call expression processing",
		}, argumentsParsingError)
	}

	var endCallToken, isEndAtEndCall = stream.Look()

	if isEndAtEndCall {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at call expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if endCallToken.Code != token.CLOSE_EXPRESSION {
		return []*ast_node.ASTNode{&callNode}, parser_error.ParserError{
			Message:       "Unknow token. Expected end call expression. But received: " + endCallToken.Code,
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	callNode.EndPosition = endCallToken.EndPosition

	if !ast_node.IsNextExpressionToken(stream) {
		return []*ast_node.ASTNode{&callNode}, nil
	}

	stream.MoveNext()

	return context.Process(stream, context, &callNode)
}

func processCallExpressionArguments(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()
	var arguments = []*ast_node.ASTNode{}

	for !isEnd && currentToken.Code != token.CLOSE_EXPRESSION {
		var argument, argumentParsingError = context.Process(stream, context, nil)

		if argumentParsingError != nil {
			return arguments, parser_error.MergeParserErrors(parser_error.ParserError{
				Message: "Something wrong at call argument processing",
			}, argumentParsingError)
		}

		if len(argument) != 1 {
			return arguments, parser_error.ParserError{
				Message:       "Parsing error. Argument declaration should have only one value node. But received: " + fmt.Sprint(len(argument)),
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}
		}

		arguments = append(arguments, argument[0])

		stream.MoveNext()
		currentToken, isEnd = stream.Look()

		if currentToken.Code != token.COMMA && currentToken.Code != token.CLOSE_EXPRESSION {
			return arguments, parser_error.ParserError{
				Message:       "Parsing error. Argument declarations should devided by comma. But received: " + currentToken.Code,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}
		}

		if currentToken.Code == token.COMMA {
			stream.MoveNext()
			currentToken, isEnd = stream.Look()
		}
	}

	return arguments, nil
}
