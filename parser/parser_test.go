package parser

import (
	"testing"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/source_mock"
)

func TestVariableDeclaration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		const a;
	`

	var parser = CreateParser(src)
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
				Params: []ast_node.ASTNodeParam{
					{
						Name: ast_node.AST_PARAM_VARIABLE_NAME,
						Value: "a",
						StartPosition: 3,
						EndPosition: 9,
					},
				},
				StartPosition: 3,
				EndPosition: 9,
			},
		},
	}

	var diff = compareAst(ast, &expectedAst)

	if len(diff) > 0 {
		t.Errorf("Different AST")
		t.Errorf("Message: %s", diff)
	}

	if err != nil {
		t.Errorf("Should parse without errors")
		t.Errorf("Failed with message: %s", err.Error())
	}
}

func compareAst(first *ast_node.ASTNode, second *ast_node.ASTNode) string {
	return compareNodes(first, second)
}

func compareNodes(first *ast_node.ASTNode, second *ast_node.ASTNode) string {
	if first.Code != second.Code {
		return "Different codes in nodes. First: " + first.Code + " Second: " + second.Code
	}

	if first.StartPosition != second.StartPosition {
		return "Different start position in nodes. First code: " + first.Code + " Second code: " + second.Code + "\n" + "First start position: " + string(first.StartPosition) + " Second: " + string(second.StartPosition) 
	}

	if first.EndPosition != second.EndPosition {
		return "Different end position in nodes. First code: " + first.Code + " Second code: " + second.Code + "\n" + "First end position: " + string(first.EndPosition) + " Second: " + string(second.EndPosition) 
	}

	var diffInParams = compareNodesParams(first, second)

	if len(diffInParams) > 0 {
		return "Diff in nodes. First: " + string(first.Code) + " Second: " + string(second.Code) + "\n" + diffInParams
	}

	if len(first.Body) != len(second.Body) {
		return "Different body size. First code: " + first.Code + " Second code: " + second.Code + "\n" + "First body size: " + string(len(first.Body)) + " Second: " + string(len(second.Body)) 
	}

	for i := 0; i < len(first.Body); i++ {
		var firstChild = first.Body[i]
		var secondChild = second.Body[i]

		var childDiff = compareNodes(firstChild, secondChild)

		if len(childDiff) > 0 {
			return "Diff in child nodes. First parent node: " + string(first.Code) + " Second parent node: " + string(second.Code) + "\n" + childDiff
		}
	}

	return ""
}

func compareNodesParams(first *ast_node.ASTNode, second *ast_node.ASTNode) string {
	if len(first.Params) != len(second.Params) {
		return "Different params count. First: " + string(len(first.Params)) + " Second: " + string(len(second.Params))
	}

	for i := 0; i < len(first.Params); i++ {
		var firstParam = first.Params[i]
		var secondParam = second.Params[i]

		if firstParam.Name != secondParam.Name {
			return "Different param name. First: " + firstParam.Name + " Second: " + secondParam.Name
		}

		if firstParam.Value != secondParam.Value {
			return "Different param value. First: " + firstParam.Value + " Second: " + secondParam.Value
		}

		if firstParam.StartPosition != secondParam.StartPosition {
			return "Different param start position. First: " + string(firstParam.StartPosition) + " Second: " + string(secondParam.StartPosition)
		}

		if firstParam.EndPosition != secondParam.EndPosition {
			return "Different param end position. First: " + string(firstParam.EndPosition) + " Second: " + string(secondParam.EndPosition)
		}
	}

	return ""
}
