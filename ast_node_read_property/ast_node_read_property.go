package ast_node_read_property

import (
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error"
)

var ReadPropertyProcessor ast_node.ASTNodeProcessor = process

func process(stream ast_node.ITokenStream, context ast_node.IASTNodeProcessingContext, leftNode *ast_node.ASTNode) ([]*ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if leftNode == nil {
		return []*ast_node.ASTNode{}, parser_error.ParserError{
			Message:       "Expected left node for read_property node",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	if isEnd {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at read property expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	var nextReadPropetryNode = ast_node.CreateNode(currentToken)
	ast_node.AppendNode(&nextReadPropetryNode, leftNode)

	stream.MoveNext()
	var propertyToken, isEndAtProperty = stream.Look()

	if isEndAtProperty {
		return []*ast_node.ASTNode{leftNode}, parser_error.ParserError{
			Message:       "Unexpected file end. Something wrong internal at read property expression processing",
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	nextReadPropetryNode.Params = []ast_node.ASTNodeParam{
		{
			Name:          ast_node.AST_PARAM_PROPERTY_NAME,
			Value:         propertyToken.Value,
			StartPosition: propertyToken.StartPosition,
			EndPosition:   propertyToken.EndPosition,
		},
	}

	if !ast_node.IsNextExpressionToken(stream) {
		return []*ast_node.ASTNode{&nextReadPropetryNode}, nil
	}

	stream.MoveNext()

	return context.Process(stream, context, &nextReadPropetryNode)
}
