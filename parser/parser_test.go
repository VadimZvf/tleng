package parser

import (
	"fmt"
	"testing"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/source_mock"
)

func TestVariableDeclaration(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	var a;
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
				Params: []ast_node.ASTNodeParam{
					{
						Name:          ast_node.AST_PARAM_VARIABLE_NAME,
						Value:         "a",
						StartPosition: 6,
						EndPosition:   6,
					},
				},
				StartPosition: 2,
				EndPosition:   6,
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

func TestNumberVariable(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	var a = 12;
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
				Params: []ast_node.ASTNodeParam{
					{
						Name:          ast_node.AST_PARAM_VARIABLE_NAME,
						Value:         "a",
						StartPosition: 6,
						EndPosition:   6,
					},
				},
				StartPosition: 2,
				EndPosition:   6,
			},
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "a",
								StartPosition: 6,
								EndPosition:   6,
							},
						},
						StartPosition: 6,
						EndPosition:   6,
					},
					{
						Code: ast_node.AST_NODE_CODE_NUMBER,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_NUMBER_VALUE,
								Value:         "12",
								StartPosition: 10,
								EndPosition:   11,
							},
						},
						StartPosition: 10,
						EndPosition:   11,
					},
				},
				StartPosition: 8,
				EndPosition:   8,
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

func TestStringVariable(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	var foo = "Hello World!";
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
				Params: []ast_node.ASTNodeParam{
					{
						Name:          ast_node.AST_PARAM_VARIABLE_NAME,
						Value:         "foo",
						StartPosition: 6,
						EndPosition:   8,
					},
				},
				StartPosition: 2,
				EndPosition:   8,
			},
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "foo",
								StartPosition: 6,
								EndPosition:   8,
							},
						},
						StartPosition: 6,
						EndPosition:   8,
					},
					{
						Code: ast_node.AST_NODE_CODE_STRING,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_STRING_VALUE,
								Value:         "Hello World!",
								StartPosition: 12,
								EndPosition:   25,
							},
						},
						StartPosition: 12,
						EndPosition:   25,
					},
				},
				StartPosition: 10,
				EndPosition:   10,
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

func TestReferenceNumberAssignment(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	bar = 777;
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "bar",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
					{
						Code: ast_node.AST_NODE_CODE_NUMBER,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_NUMBER_VALUE,
								Value:         "777",
								StartPosition: 8,
								EndPosition:   10,
							},
						},
						StartPosition: 8,
						EndPosition:   10,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestNumberSumm(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	bar = 3 + 9;
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "bar",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
					{
						Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
								Value:         "+",
								StartPosition: 10,
								EndPosition:   10,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "3",
										StartPosition: 8,
										EndPosition:   8,
									},
								},
								StartPosition: 8,
								EndPosition:   8,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "9",
										StartPosition: 12,
										EndPosition:   12,
									},
								},
								StartPosition: 12,
								EndPosition:   12,
							},
						},
						StartPosition: 10,
						EndPosition:   10,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestReferenceSumm(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	bar = foo + baz;
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "bar",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
					{
						Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
								Value:         "+",
								StartPosition: 12,
								EndPosition:   12,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "foo",
										StartPosition: 8,
										EndPosition:   10,
									},
								},
								StartPosition: 8,
								EndPosition:   10,
							},
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "baz",
										StartPosition: 14,
										EndPosition:   16,
									},
								},
								StartPosition: 14,
								EndPosition:   16,
							},
						},
						StartPosition: 12,
						EndPosition:   12,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestReferenceWithNumberSumm(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	bar = foo + 55;
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "bar",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
					{
						Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
								Value:         "+",
								StartPosition: 12,
								EndPosition:   12,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "foo",
										StartPosition: 8,
										EndPosition:   10,
									},
								},
								StartPosition: 8,
								EndPosition:   10,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "55",
										StartPosition: 14,
										EndPosition:   15,
									},
								},
								StartPosition: 14,
								EndPosition:   15,
							},
						},
						StartPosition: 12,
						EndPosition:   12,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestParenthesizedExpression(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	bar = (foo + 55);
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "bar",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
					{
						Code: ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION,
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
										Value:         "+",
										StartPosition: 13,
										EndPosition:   13,
									},
								},
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_REFERENCE,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_VARIABLE_NAME,
												Value:         "foo",
												StartPosition: 9,
												EndPosition:   11,
											},
										},
										StartPosition: 9,
										EndPosition:   11,
									},
									{
										Code: ast_node.AST_NODE_CODE_NUMBER,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_NUMBER_VALUE,
												Value:         "55",
												StartPosition: 15,
												EndPosition:   16,
											},
										},
										StartPosition: 15,
										EndPosition:   16,
									},
								},
								StartPosition: 13,
								EndPosition:   13,
							},
						},
						StartPosition: 8,
						EndPosition:   17,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestTwoParenthesizedExpression(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	baz = (foo + 55) / (9 - 3);
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "baz",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
					{
						Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
								Value:         "/",
								StartPosition: 19,
								EndPosition:   19,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION,
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
												Value:         "+",
												StartPosition: 13,
												EndPosition:   13,
											},
										},
										Body: []*ast_node.ASTNode{
											{
												Code: ast_node.AST_NODE_CODE_REFERENCE,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_VARIABLE_NAME,
														Value:         "foo",
														StartPosition: 9,
														EndPosition:   11,
													},
												},
												StartPosition: 9,
												EndPosition:   11,
											},
											{
												Code: ast_node.AST_NODE_CODE_NUMBER,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_NUMBER_VALUE,
														Value:         "55",
														StartPosition: 15,
														EndPosition:   16,
													},
												},
												StartPosition: 15,
												EndPosition:   16,
											},
										},
										StartPosition: 13,
										EndPosition:   13,
									},
								},
								StartPosition: 8,
								EndPosition:   17,
							},
							{
								Code: ast_node.AST_NODE_CODE_PARENTHESIZED_EXPRESSION,
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
												Value:         "-",
												StartPosition: 24,
												EndPosition:   24,
											},
										},
										Body: []*ast_node.ASTNode{
											{
												Code: ast_node.AST_NODE_CODE_NUMBER,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_NUMBER_VALUE,
														Value:         "9",
														StartPosition: 22,
														EndPosition:   22,
													},
												},
												StartPosition: 22,
												EndPosition:   22,
											},
											{
												Code: ast_node.AST_NODE_CODE_NUMBER,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_NUMBER_VALUE,
														Value:         "3",
														StartPosition: 26,
														EndPosition:   26,
													},
												},
												StartPosition: 26,
												EndPosition:   26,
											},
										},
										StartPosition: 24,
										EndPosition:   24,
									},
								},
								StartPosition: 21,
								EndPosition:   27,
							},
						},
						StartPosition: 19,
						EndPosition:   19,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestFunctionDeclaration(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	function baz() {
		var a = 1;
		b = 2;
	};
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_FUNCTION,
				Params: []ast_node.ASTNodeParam{
					{
						Name:          ast_node.AST_PARAM_FUNCTION_NAME,
						Value:         "baz",
						StartPosition: 11,
						EndPosition:   13,
					},
				},
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "a",
								StartPosition: 25,
								EndPosition:   25,
							},
						},
						StartPosition: 21,
						EndPosition:   25,
					},
					{
						Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "a",
										StartPosition: 25,
										EndPosition:   25,
									},
								},
								StartPosition: 25,
								EndPosition:   25,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "1",
										StartPosition: 29,
										EndPosition:   29,
									},
								},
								StartPosition: 29,
								EndPosition:   29,
							},
						},
						StartPosition: 27,
						EndPosition:   27,
					},
					{
						Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "b",
										StartPosition: 34,
										EndPosition:   34,
									},
								},
								StartPosition: 34,
								EndPosition:   34,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "2",
										StartPosition: 38,
										EndPosition:   38,
									},
								},
								StartPosition: 38,
								EndPosition:   38,
							},
						},
						StartPosition: 36,
						EndPosition:   36,
					},
				},
				StartPosition: 2,
				EndPosition:   42,
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

func TestFunctionDeclarationWithReturn(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	function baz() {
		return "Hello" + "World"
	};
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_FUNCTION,
				Params: []ast_node.ASTNodeParam{
					{
						Name:          ast_node.AST_PARAM_FUNCTION_NAME,
						Value:         "baz",
						StartPosition: 11,
						EndPosition:   13,
					},
				},
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_RETURN,
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
										Value:         "+",
										StartPosition: 36,
										EndPosition:   36,
									},
								},
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_STRING,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_STRING_VALUE,
												Value:         "Hello",
												StartPosition: 28,
												EndPosition:   34,
											},
										},
										StartPosition: 28,
										EndPosition:   34,
									},
									{
										Code: ast_node.AST_NODE_CODE_STRING,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_STRING_VALUE,
												Value:         "World",
												StartPosition: 38,
												EndPosition:   44,
											},
										},
										StartPosition: 38,
										EndPosition:   44,
									},
								},
								StartPosition: 36,
								EndPosition:   36,
							},
						},
						StartPosition: 21,
						EndPosition:   26,
					},
				},
				StartPosition: 2,
				EndPosition:   47,
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

func TestReadProperty(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	a.b + 23
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
				Params: []ast_node.ASTNodeParam{
					{
						Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
						Value:         "+",
						StartPosition: 6,
						EndPosition:   6,
					},
				},
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_READ_PROP,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_PROPERTY_NAME,
								Value:         "b",
								StartPosition: 4,
								EndPosition:   4,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "a",
										StartPosition: 2,
										EndPosition:   2,
									},
								},
								StartPosition: 2,
								EndPosition:   2,
							},
						},
						StartPosition: 3,
						EndPosition:   3,
					},
					{
						Code: ast_node.AST_NODE_CODE_NUMBER,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_NUMBER_VALUE,
								Value:         "23",
								StartPosition: 8,
								EndPosition:   9,
							},
						},
						StartPosition: 8,
						EndPosition:   9,
					},
				},
				StartPosition: 6,
				EndPosition:   6,
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

func TestReadTwoProperties(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	a.foo.bar = 12
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_ASSIGNMENT,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_READ_PROP,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_PROPERTY_NAME,
								Value:         "bar",
								StartPosition: 8,
								EndPosition:   10,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_READ_PROP,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_PROPERTY_NAME,
										Value:         "foo",
										StartPosition: 4,
										EndPosition:   6,
									},
								},
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_REFERENCE,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_VARIABLE_NAME,
												Value:         "a",
												StartPosition: 2,
												EndPosition:   2,
											},
										},
										StartPosition: 2,
										EndPosition:   2,
									},
								},
								StartPosition: 3,
								EndPosition:   3,
							},
						},
						StartPosition: 7,
						EndPosition:   7,
					},
					{
						Code: ast_node.AST_NODE_CODE_NUMBER,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_NUMBER_VALUE,
								Value:         "12",
								StartPosition: 14,
								EndPosition:   15,
							},
						},
						StartPosition: 14,
						EndPosition:   15,
					},
				},
				StartPosition: 12,
				EndPosition:   12,
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

func TestCallFunction(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	a()
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_CALL_EXPRESSION,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "a",
								StartPosition: 2,
								EndPosition:   2,
							},
						},
						StartPosition: 2,
						EndPosition:   2,
					},
				},
				StartPosition: 3,
				EndPosition:   4,
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

func TestCallFunctionWithProperty(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	bar.baz()
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_CALL_EXPRESSION,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_READ_PROP,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_PROPERTY_NAME,
								Value:         "baz",
								StartPosition: 6,
								EndPosition:   8,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "bar",
										StartPosition: 2,
										EndPosition:   4,
									},
								},
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 5,
						EndPosition:   5,
					},
				},
				StartPosition: 9,
				EndPosition:   10,
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

func TestCallFunctionWithArguments(t *testing.T) {
	var src = source_mock.GetSourceMock(`
	baz(bar(), "foo")
`)
	var parser = CreateParser(src, createMockStdout())
	var ast, err = parser.Parse(false)

	var expectedAst = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
		Body: []*ast_node.ASTNode{
			{
				Code: ast_node.AST_NODE_CODE_CALL_EXPRESSION,
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_REFERENCE,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "baz",
								StartPosition: 2,
								EndPosition:   4,
							},
						},
						StartPosition: 2,
						EndPosition:   4,
					},
				},
				Arguments: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_CALL_EXPRESSION,
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "bar",
										StartPosition: 6,
										EndPosition:   8,
									},
								},
								StartPosition: 6,
								EndPosition:   8,
							},
						},
						StartPosition: 9,
						EndPosition:   10,
					},
					{
						Code: ast_node.AST_NODE_CODE_STRING,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_STRING_VALUE,
								Value:         "foo",
								StartPosition: 13,
								EndPosition:   17,
							},
						},
						StartPosition: 13,
						EndPosition:   17,
					},
				},
				StartPosition: 5,
				EndPosition:   18,
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
		return "Different start position in nodes. First code: " + first.Code + " Second code: " + second.Code + "\n" + "First start position: " + fmt.Sprint(first.StartPosition) + " Second: " + fmt.Sprint(second.StartPosition)
	}

	if first.EndPosition != second.EndPosition {
		return "Different end position in nodes. First code: " + first.Code + " Second code: " + second.Code + "\n" + "First end position: " + fmt.Sprint(first.EndPosition) + " Second: " + fmt.Sprint(second.EndPosition)
	}

	var diffInParams = compareNodesParams(first, second)

	if len(diffInParams) > 0 {
		return "Diff in nodes. First: " + first.Code + " Second: " + second.Code + "\n" + diffInParams
	}

	if len(first.Body) != len(second.Body) {
		return "Different body size. First code: " + first.Code + " Second code: " + second.Code + "\n" + "First body size: " + fmt.Sprint(len(first.Body)) + " Second: " + fmt.Sprint(len(second.Body))
	}

	for i := 0; i < len(first.Body); i++ {
		var firstChild = first.Body[i]
		var secondChild = second.Body[i]

		var childDiff = compareNodes(firstChild, secondChild)

		if len(childDiff) > 0 {
			return "Diff in child nodes. First parent node: " + first.Code + " Second parent node: " + second.Code + "\n" + childDiff
		}
	}

	return ""
}

func compareNodesParams(first *ast_node.ASTNode, second *ast_node.ASTNode) string {
	if len(first.Params) != len(second.Params) {
		return "Different params count. First: " + fmt.Sprint(len(first.Params)) + " Second: " + fmt.Sprint(len(second.Params))
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
			return "Different param start position. First: " + fmt.Sprint(firstParam.StartPosition) + " Second: " + fmt.Sprint(secondParam.StartPosition)
		}

		if firstParam.EndPosition != secondParam.EndPosition {
			return "Different param end position. First: " + fmt.Sprint(firstParam.EndPosition) + " Second: " + fmt.Sprint(secondParam.EndPosition)
		}
	}

	return ""
}

type mockStdout struct{}

func (std *mockStdout) Print(line string)      {}
func (std *mockStdout) PrintError(line string) {}

func createMockStdout() *mockStdout {
	return &mockStdout{}
}
