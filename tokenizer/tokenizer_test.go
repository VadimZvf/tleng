package tokenizer

import (
	"testing"

	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
	"github.com/VadimZvf/golang/token_variable_decloration"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestStringVariable(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		const a = "Hello world";
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token.NEW_LINE,
		Position: 0,
	})

	compareTokens(t, tokens[1], token.Token{
		Code:     token_variable_decloration.VARIABLE_DECLARAION,
		Position: 7,
	})

	checkParam(t, tokens[1], token.TokenParam{
		Name:     token_variable_decloration.VARIABLE_NAME_PARAM,
		Value:    "a",
		Position: 9,
	})

	compareTokens(t, tokens[2], token.Token{
		Code:     token.ASSIGNMENT,
		Position: 11,
	})

	compareTokens(t, tokens[3], token.Token{
		Code:     token_string.STRING,
		Position: 25,
		Value:    "Hello world",
	})
}

func TestKeywordLikeVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `constA;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token.KEY_WORD,
		Value:    "constA",
		Position: 5,
	})
}

func TestKeywordLikeFunctionDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `functionA();`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token.KEY_WORD,
		Value:    "functionA",
		Position: 8,
	})

	compareTokens(t, tokens[1], token.Token{
		Code:     token.OPEN_EXPRESSION,
		Position: 9,
	})

	compareTokens(t, tokens[2], token.Token{
		Code:     token.CLOSE_EXPRESSION,
		Position: 10,
	})
}

func TestCopyVariableByLink(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const a = b;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token_variable_decloration.VARIABLE_DECLARAION,
		Position: 4,
	})

	checkParam(t, tokens[0], token.TokenParam{
		Name:     token_variable_decloration.VARIABLE_NAME_PARAM,
		Value:    "a",
		Position: 6,
	})

	compareTokens(t, tokens[1], token.Token{
		Code:     token.ASSIGNMENT,
		Position: 8,
	})

	compareTokens(t, tokens[2], token.Token{
		Code:     token.KEY_WORD,
		Position: 10,
		Value:    "b",
	})
}

func TestVariableWithTokenName(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const functionA;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token_variable_decloration.VARIABLE_DECLARAION,
		Position: 4,
	})

	checkParam(t, tokens[0], token.TokenParam{
		Name:     token_variable_decloration.VARIABLE_NAME_PARAM,
		Value:    "functionA",
		Position: 14,
	})
}

func TestFunctionDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo() {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token_function_declaration.FUNCTION_DECLARATION,
		Position: 7,
	})

	checkParam(t, tokens[0], token.TokenParam{
		Name:     token_function_declaration.FUNCTION_NAME_PARAM,
		Value:    "foo",
		Position: 11,
	})

	compareTokens(t, tokens[1], token.Token{
		Code:     token.OPEN_BLOCK,
		Position: 15,
	})

	compareTokens(t, tokens[2], token.Token{
		Code:     token.CLOSE_BLOCK,
		Position: 16,
	})
}

func TestFunctionDeclorationWithArguments(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo(a, b) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token_function_declaration.FUNCTION_DECLARATION,
		Position: 7,
	})

	checkParam(t, tokens[0], token.TokenParam{
		Name:     token_function_declaration.FUNCTION_NAME_PARAM,
		Value:    "foo",
		Position: 11,
	})

	checkParam(t, tokens[0], token.TokenParam{
		Name:     token_function_declaration.FUNCTION_ARGUMENT_PARAM,
		Value:    "a",
		Position: 13,
	})

	checkParam(t, tokens[0], token.TokenParam{
		Name:     token_function_declaration.FUNCTION_ARGUMENT_PARAM,
		Value:    "b",
		Position: 16,
	})

	compareTokens(t, tokens[1], token.Token{
		Code:     token.OPEN_BLOCK,
		Position: 19,
	})

	compareTokens(t, tokens[2], token.Token{
		Code:     token.CLOSE_BLOCK,
		Position: 20,
	})
}

func TestFunctionWithReturnStatement(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		function bar() {
			return "Avada kedavra"
		}
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	compareTokens(t, tokens[0], token.Token{
		Code:     token.NEW_LINE,
		Position: 0,
	})

	compareTokens(t, tokens[1], token.Token{
		Code:     token_function_declaration.FUNCTION_DECLARATION,
		Position: 10,
	})

	checkParam(t, tokens[1], token.TokenParam{
		Name:     token_function_declaration.FUNCTION_NAME_PARAM,
		Value:    "bar",
		Position: 14,
	})

	compareTokens(t, tokens[2], token.Token{
		Code:     token.OPEN_BLOCK,
		Position: 18,
	})

	compareTokens(t, tokens[3], token.Token{
		Code:     token.NEW_LINE,
		Position: 19,
	})

	compareTokens(t, tokens[4], token.Token{
		Code:     token_return.RETURN_DECLARATION,
		Position: 28,
	})

	compareTokens(t, tokens[5], token.Token{
		Code:     token_string.STRING,
		Position: 44,
		Value:    "Avada kedavra",
	})

	compareTokens(t, tokens[6], token.Token{
		Code:     token.NEW_LINE,
		Position: 45,
	})

	compareTokens(t, tokens[7], token.Token{
		Code:     token.CLOSE_BLOCK,
		Position: 48,
	})
}

/// Utils

func compareTokens(t *testing.T, first token.Token, second token.Token) {
	if first.Code != second.Code {
		t.Errorf("Different token Codes: %s - %s", first.Code, second.Code)
	}

	if first.Value != second.Value {
		t.Errorf("Different token Values: %s - %s", first.Value, second.Value)
	}

	if first.Position != second.Position {
		t.Errorf("Different token position: %d - %d", first.Position, second.Position)
	}
}

func checkParam(t *testing.T, tokenItem token.Token, targetParam token.TokenParam) {
	for _, param := range tokenItem.Params {
		if param.Name == targetParam.Name && param.Value == targetParam.Value && param.Position == targetParam.Position {
			return
		}
	}

	t.Errorf("Param not found in token")
	t.Errorf("Looking for param - name: %s, value: %s, position, %d", targetParam.Name, targetParam.Value, targetParam.Position)
	t.Errorf("Token has params:")

	for _, param := range tokenItem.Params {
		t.Errorf("name: %s, value: %s, position, %d", param.Name, param.Value, param.Position)
	}
}
