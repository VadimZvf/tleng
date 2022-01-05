package parser

import (
	"fmt"
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
						Name:          ast_node.AST_PARAM_VARIABLE_NAME,
						Value:         "a",
						StartPosition: 9,
						EndPosition:   9,
					},
				},
				StartPosition: 3,
				EndPosition:   9,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		const a = 12;
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
						Name:          ast_node.AST_PARAM_VARIABLE_NAME,
						Value:         "a",
						StartPosition: 9,
						EndPosition:   9,
					},
				},
				StartPosition: 3,
				EndPosition:   9,
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
								StartPosition: 9,
								EndPosition:   9,
							},
						},
						StartPosition: 9,
						EndPosition:   9,
					},
					{
						Code: ast_node.AST_NODE_CODE_NUMBER,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_NUMBER_VALUE,
								Value:         "12",
								StartPosition: 13,
								EndPosition:   14,
							},
						},
						StartPosition: 13,
						EndPosition:   14,
					},
				},
				StartPosition: 11,
				EndPosition:   11,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		const foo = "Hello World!";
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
						Name:          ast_node.AST_PARAM_VARIABLE_NAME,
						Value:         "foo",
						StartPosition: 9,
						EndPosition:   11,
					},
				},
				StartPosition: 3,
				EndPosition:   11,
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
								StartPosition: 9,
								EndPosition:   11,
							},
						},
						StartPosition: 9,
						EndPosition:   11,
					},
					{
						Code: ast_node.AST_NODE_CODE_STRING,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_STRING_VALUE,
								Value:         "Hello World!",
								StartPosition: 15,
								EndPosition:   28,
							},
						},
						StartPosition: 15,
						EndPosition:   28,
					},
				},
				StartPosition: 13,
				EndPosition:   13,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		bar = 777;
	`

	var parser = CreateParser(src)
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
								StartPosition: 3,
								EndPosition:   5,
							},
						},
						StartPosition: 3,
						EndPosition:   5,
					},
					{
						Code: ast_node.AST_NODE_CODE_NUMBER,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_NUMBER_VALUE,
								Value:         "777",
								StartPosition: 9,
								EndPosition:   11,
							},
						},
						StartPosition: 9,
						EndPosition:   11,
					},
				},
				StartPosition: 7,
				EndPosition:   7,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		bar = 3 + 9;
	`

	var parser = CreateParser(src)
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
								StartPosition: 3,
								EndPosition:   5,
							},
						},
						StartPosition: 3,
						EndPosition:   5,
					},
					{
						Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
								Value:         "+",
								StartPosition: 11,
								EndPosition:   11,
							},
						},
						Body: []*ast_node.ASTNode{
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "3",
										StartPosition: 9,
										EndPosition:   9,
									},
								},
								StartPosition: 9,
								EndPosition:   9,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "9",
										StartPosition: 13,
										EndPosition:   13,
									},
								},
								StartPosition: 13,
								EndPosition:   13,
							},
						},
						StartPosition: 11,
						EndPosition:   11,
					},
				},
				StartPosition: 7,
				EndPosition:   7,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		bar = foo + baz;
	`

	var parser = CreateParser(src)
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
								StartPosition: 3,
								EndPosition:   5,
							},
						},
						StartPosition: 3,
						EndPosition:   5,
					},
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
								Code: ast_node.AST_NODE_CODE_REFERENCE,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_VARIABLE_NAME,
										Value:         "baz",
										StartPosition: 15,
										EndPosition:   17,
									},
								},
								StartPosition: 15,
								EndPosition:   17,
							},
						},
						StartPosition: 13,
						EndPosition:   13,
					},
				},
				StartPosition: 7,
				EndPosition:   7,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		bar = foo + 55;
	`

	var parser = CreateParser(src)
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
								StartPosition: 3,
								EndPosition:   5,
							},
						},
						StartPosition: 3,
						EndPosition:   5,
					},
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
				StartPosition: 7,
				EndPosition:   7,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		bar = (foo + 55);
	`

	var parser = CreateParser(src)
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
								StartPosition: 3,
								EndPosition:   5,
							},
						},
						StartPosition: 3,
						EndPosition:   5,
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
										StartPosition: 14,
										EndPosition:   14,
									},
								},
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_REFERENCE,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_VARIABLE_NAME,
												Value:         "foo",
												StartPosition: 10,
												EndPosition:   12,
											},
										},
										StartPosition: 10,
										EndPosition:   12,
									},
									{
										Code: ast_node.AST_NODE_CODE_NUMBER,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_NUMBER_VALUE,
												Value:         "55",
												StartPosition: 16,
												EndPosition:   17,
											},
										},
										StartPosition: 16,
										EndPosition:   17,
									},
								},
								StartPosition: 14,
								EndPosition:   14,
							},
						},
						StartPosition: 9,
						EndPosition:   18,
					},
				},
				StartPosition: 7,
				EndPosition:   7,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		baz = (foo + 55) / (9 - 3);
	`

	var parser = CreateParser(src)
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
								StartPosition: 3,
								EndPosition:   5,
							},
						},
						StartPosition: 3,
						EndPosition:   5,
					},
					{
						Code: ast_node.AST_NODE_CODE_BINARY_EXPRESSION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_BINARY_EXPRESSION_TYPE,
								Value:         "/",
								StartPosition: 20,
								EndPosition:   20,
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
												StartPosition: 14,
												EndPosition:   14,
											},
										},
										Body: []*ast_node.ASTNode{
											{
												Code: ast_node.AST_NODE_CODE_REFERENCE,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_VARIABLE_NAME,
														Value:         "foo",
														StartPosition: 10,
														EndPosition:   12,
													},
												},
												StartPosition: 10,
												EndPosition:   12,
											},
											{
												Code: ast_node.AST_NODE_CODE_NUMBER,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_NUMBER_VALUE,
														Value:         "55",
														StartPosition: 16,
														EndPosition:   17,
													},
												},
												StartPosition: 16,
												EndPosition:   17,
											},
										},
										StartPosition: 14,
										EndPosition:   14,
									},
								},
								StartPosition: 9,
								EndPosition:   18,
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
												StartPosition: 25,
												EndPosition:   25,
											},
										},
										Body: []*ast_node.ASTNode{
											{
												Code: ast_node.AST_NODE_CODE_NUMBER,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_NUMBER_VALUE,
														Value:         "9",
														StartPosition: 23,
														EndPosition:   23,
													},
												},
												StartPosition: 23,
												EndPosition:   23,
											},
											{
												Code: ast_node.AST_NODE_CODE_NUMBER,
												Params: []ast_node.ASTNodeParam{
													{
														Name:          ast_node.AST_PARAM_NUMBER_VALUE,
														Value:         "3",
														StartPosition: 27,
														EndPosition:   27,
													},
												},
												StartPosition: 27,
												EndPosition:   27,
											},
										},
										StartPosition: 25,
										EndPosition:   25,
									},
								},
								StartPosition: 22,
								EndPosition:   28,
							},
						},
						StartPosition: 20,
						EndPosition:   20,
					},
				},
				StartPosition: 7,
				EndPosition:   7,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		function baz() {
			const a = 1;
			b = 2;
		};
	`

	var parser = CreateParser(src)
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
						StartPosition: 12,
						EndPosition:   14,
					},
				},
				Body: []*ast_node.ASTNode{
					{
						Code: ast_node.AST_NODE_CODE_VARIABLE_DECLARATION,
						Params: []ast_node.ASTNodeParam{
							{
								Name:          ast_node.AST_PARAM_VARIABLE_NAME,
								Value:         "a",
								StartPosition: 29,
								EndPosition:   29,
							},
						},
						StartPosition: 23,
						EndPosition:   29,
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
										StartPosition: 29,
										EndPosition:   29,
									},
								},
								StartPosition: 29,
								EndPosition:   29,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "1",
										StartPosition: 33,
										EndPosition:   33,
									},
								},
								StartPosition: 33,
								EndPosition:   33,
							},
						},
						StartPosition: 31,
						EndPosition:   31,
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
										StartPosition: 39,
										EndPosition:   39,
									},
								},
								StartPosition: 39,
								EndPosition:   39,
							},
							{
								Code: ast_node.AST_NODE_CODE_NUMBER,
								Params: []ast_node.ASTNodeParam{
									{
										Name:          ast_node.AST_PARAM_NUMBER_VALUE,
										Value:         "2",
										StartPosition: 43,
										EndPosition:   43,
									},
								},
								StartPosition: 43,
								EndPosition:   43,
							},
						},
						StartPosition: 41,
						EndPosition:   41,
					},
				},
				StartPosition: 3,
				EndPosition:   48,
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
	var src = source_mock.GetSourceMock()
	src.FullText = `
		function baz() {
			return "Hello" + "World"
		};
	`

	var parser = CreateParser(src)
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
						StartPosition: 12,
						EndPosition:   14,
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
										StartPosition: 38,
										EndPosition:   38,
									},
								},
								Body: []*ast_node.ASTNode{
									{
										Code: ast_node.AST_NODE_CODE_STRING,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_STRING_VALUE,
												Value:         "Hello",
												StartPosition: 30,
												EndPosition:   36,
											},
										},
										StartPosition: 30,
										EndPosition:   36,
									},
									{
										Code: ast_node.AST_NODE_CODE_STRING,
										Params: []ast_node.ASTNodeParam{
											{
												Name:          ast_node.AST_PARAM_STRING_VALUE,
												Value:         "World",
												StartPosition: 40,
												EndPosition:   46,
											},
										},
										StartPosition: 40,
										EndPosition:   46,
									},
								},
								StartPosition: 38,
								EndPosition:   38,
							},
						},
						StartPosition: 23,
						EndPosition:   28,
					},
				},
				StartPosition: 3,
				EndPosition:   50,
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
