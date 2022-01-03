package ast

import (
	"github.com/VadimZvf/golang/ast_error"
	"github.com/VadimZvf/golang/ast_factory"
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/ast_token_stream"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_keyword"
	"github.com/VadimZvf/golang/token_number"
	"github.com/VadimZvf/golang/token_read_property"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
	"github.com/VadimZvf/golang/token_variable_declaration"
)

func CreateAST(tokens []token.Token) (*ast_node.ASTNode, error) {
	var tokenStream = ast_token_stream.CreateTokenStream(tokens)
	var _, isEnd = tokenStream.Look()
	var factory = ast_factory.CreateASTFactory()

	for !isEnd {
		var node, err = getNode(&tokenStream)

		factory.Append(&node)

		if err != nil {
			return factory.GetAST(), err
		}

		tokenStream.MoveNext()
		_, isEnd = tokenStream.Look()
	}

	return factory.GetAST(), nil
}

func getNode(stream *ast_token_stream.TokenStream) (ast_node.ASTNode, error) {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}, ast_error.AstError{
			Message: "File ended",
		}
	}


	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return ast_node.ASTNode{}, ast_error.AstError{
			Message: "File ended",
		}
	}

	switch nextToken.Code {
	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		var leftNode = createNode(currentToken)
		stream.MoveNext()
		var binaryNode = createNode(nextToken)
		stream.MoveNext()

		var rightNode, err = getNode(stream)
	
		if err != nil {
			return ast_node.ASTNode{}, err
		}

		appendNodes(&binaryNode, []*ast_node.ASTNode{&leftNode, &rightNode})

		return binaryNode, nil

	case token_read_property.READ_PROPERTY:
		var leftNode = createNode(currentToken)
		stream.MoveNext()
		var readNode = createNode(nextToken)
		stream.MoveNext()

		var rightNodes, err = getNode(stream)
	
		if err != nil {
			return ast_node.ASTNode{}, err
		}

		appendNodes(&readNode, []*ast_node.ASTNode{&leftNode, &rightNodes})

		return readNode, nil
	}

	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		var variableNode = createNode(currentToken)

		var nextToken, isEnd = stream.LookNext()

		if isEnd {
			break
		}

		if nextToken.Code == token.ASSIGNMENT {
			stream.MoveNext()

			var variableNameParam = token_variable_declaration.GetVariableNameParam(currentToken)
			var assignmentNode = createNode(nextToken)
			assignmentNode.Params = []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         variableNameParam.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}}

			stream.MoveNext()
			var valueNode, err = getNode(stream)

			if err != nil {
				return valueNode, err
			}

			assignmentNode.Body = []*ast_node.ASTNode{&valueNode}

			stream.MoveNext()
			return assignmentNode, nil
		} else {
			return variableNode, nil
		}

	case token_number.NUMBER:
		var numberNode = createNode(currentToken)
		return numberNode, nil

	case token_string.STRING:
		var stringNode = createNode(currentToken)
		return stringNode, nil		

	case token_return.RETURN_DECLARATION:
		var returnNode = createNode(currentToken)
		stream.MoveNext()
		var returnValueNode, err = getNode(stream)

		if err != nil {
			return returnNode, err
		}

		returnNode.Body = []*ast_node.ASTNode{&returnValueNode}
		return returnNode, nil

	case token_keyword.KEY_WORD:
		var keyWordNode = processKeyWordToken(stream)
		return keyWordNode, nil

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionNode = processFunctionToken(stream)
		return functionNode, nil

	case token.END_LINE:
		// Skip node
		// stream.MoveNext()

	}

	return ast_node.ASTNode{}, nil
}

func processKeyWordToken(stream *ast_token_stream.TokenStream) ast_node.ASTNode {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}
	}

	var nextToken, isEndNext = stream.LookNext()

	if isEndNext {
		return createNode(currentToken)
	}

	switch nextToken.Code { 
	case token.ASSIGNMENT:
		var assignmentNode = createNode(nextToken)
		assignmentNode.Params = []ast_node.ASTNodeParam{{
			Name:          ast_node.AST_PARAM_VARIABLE_NAME,
			Value:         currentToken.Value,
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}}

		stream.MoveNext()
		stream.MoveNext()

		var valueNode, error = getNode(stream);

		if error != nil {
			return ast_node.ASTNode{}
		}

		assignmentNode.Body = []*ast_node.ASTNode{&valueNode}

		return assignmentNode

	case token.OPEN_EXPRESSION:
		stream.MoveNext()
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_FUNCTION_CALL,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_FUNCTION_NAME,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	}

	return createNode(currentToken)
}

func processFunctionToken(stream *ast_token_stream.TokenStream) ast_node.ASTNode {
	var currentToken, isEnd = stream.Look()

	if isEnd {
		return ast_node.ASTNode{}
	}

	var functionNode = createNode(currentToken)

	stream.MoveNext()
	currentToken, isEnd = stream.Look()

	for !isEnd && currentToken.Code != token.CLOSE_BLOCK {
		stream.MoveNext()
		currentToken, isEnd = stream.Look()

		var node, err = getNode(stream)

		if err != nil {
			return functionNode
		}

		appendNode(&functionNode, &node)
	}

	return functionNode
}

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
// ┃         Utilities          ┃
// ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

func appendNodes(node *ast_node.ASTNode, children []*ast_node.ASTNode)  {
	for _, child := range children {
		node.Body = append(node.Body, child)
	}
}

func appendNode(node *ast_node.ASTNode, child *ast_node.ASTNode)  {
	node.Body = append(node.Body, child)
}


func createNode(currentToken token.Token) ast_node.ASTNode {
	switch currentToken.Code {
	case token_variable_declaration.VARIABLE_DECLARAION:
		var variableName = token_variable_declaration.GetVariableNameParam(currentToken)
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         variableName.Value,
				StartPosition: variableName.StartPosition,
				EndPosition:   variableName.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.ASSIGNMENT:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token.ADD, token.SUBTRACT, token.SLASH, token.ASTERISK:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_read_property.READ_PROPERTY:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_READ_PROP,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_number.NUMBER:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_NUMBER,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_NUMBER_VALUE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_string.STRING:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_STRING,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_STRING_VALUE,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_keyword.KEY_WORD:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_REFERENCE,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_VARIABLE_NAME,
				Value:         currentToken.Value,
				StartPosition: currentToken.StartPosition,
				EndPosition:   currentToken.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_return.RETURN_DECLARATION:
		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_RETURN,
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}

	case token_function_declaration.FUNCTION_DECLARATION:
		var functionName = token_function_declaration.GetFunctionNameParam(currentToken)

		return ast_node.ASTNode{
			Code: ast_node.AST_NODE_CODE_FUNCTION,
			Params: []ast_node.ASTNodeParam{{
				Name:          ast_node.AST_PARAM_FUNCTION_NAME,
				Value:         functionName.Value,
				StartPosition: functionName.StartPosition,
				EndPosition:   functionName.EndPosition,
			}},
			// Debug data
			StartPosition: currentToken.StartPosition,
			EndPosition:   currentToken.EndPosition,
		}
	default:
		return ast_node.ASTNode{}
	}
}
