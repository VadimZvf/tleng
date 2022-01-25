package ast_node_variable_declaration

import (
	"fmt"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var VariableDeclarationProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at variable declaration processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var variableDeclarationNode = ast_node.CreateNode(currentToken)

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext || nextToken.Code != token.ASSIGNMENT {
		return []*ast_node.ASTNode{&variableDeclarationNode}, nil
	}

	var variableNameParam = ast_node.GetVariableNameParam(&variableDeclarationNode)
	var referenceNode = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_REFERENCE,
		Params: []ast_node.ASTNodeParam{{
			Name:          ast_node.AST_PARAM_VARIABLE_NAME,
			Value:         variableNameParam.Value,
			StartPosition: variableNameParam.StartPosition,
			EndPosition:   variableNameParam.EndPosition,
		}},
		StartPosition: variableNameParam.StartPosition,
		EndPosition:   variableNameParam.EndPosition,
	}

	stream.MoveNext()

	var assignmentNodes, assignmentNodeParsingError = context.Process(stream, context, &referenceNode)

	if assignmentNodeParsingError != nil {
		return []*ast_node.ASTNode{&variableDeclarationNode}, parser_error.MergeParserErrors(parser_error.ParserError{
			Message: "Parsing error. At assignment with variable declaration",
		}, assignmentNodeParsingError)
	}

	if len(assignmentNodes) != 1 {
		return []*ast_node.ASTNode{&variableDeclarationNode}, parser_error.ParserError{
			Message:       "Parsing error. Should assign only one node. But received: " + fmt.Sprint(len(assignmentNodes)),
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	return []*ast_node.ASTNode{&variableDeclarationNode, assignmentNodes[0]}, nil
}
