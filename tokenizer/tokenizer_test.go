package tokenizer

import (
	"fmt"
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
	"github.com/VadimZvf/golang/token_variable_declaration"
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

	if !isSameToken(tokens[0], token.Token{
		Code:          token_variable_declaration.VARIABLE_DECLARAION,
		StartPosition: 3,
		EndPosition:   9,
	}) {
		t.Errorf("Wrong token")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_variable_declaration.VARIABLE_NAME_PARAM,
		Value:         "a",
		StartPosition: 9,
		EndPosition:   9,
	}) {
		t.Errorf("Wrong param")
	}

	if !isSameToken(tokens[1], token.Token{
		Code:          token.ASSIGNMENT,
		Value: "=",
		StartPosition: 11,
		EndPosition:   11,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[2], token.Token{
		Code:          token_string.STRING,
		Value:         "Hello world",
		StartPosition: 13,
		EndPosition:   25,
	}) {
		t.Errorf("Wrong token")
	}
}

func TestKeywordLikeVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `constA;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	if !isSameToken(tokens[0], token.Token{
		Code:          token.KEY_WORD,
		Value:         "constA",
		StartPosition: 0,
		EndPosition:   5,
	}) {
		t.Errorf("Wrong token")
	}
}

func TestKeywordLikeFunctionDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `functionA();`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	if !isSameToken(tokens[0], token.Token{
		Code:          token.KEY_WORD,
		Value:         "functionA",
		StartPosition: 0,
		EndPosition:   8,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[1], token.Token{
		Code:          token.OPEN_EXPRESSION,
		Value: "(",
		StartPosition: 9,
		EndPosition:   9,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[2], token.Token{
		Code:          token.CLOSE_EXPRESSION,
		Value: ")",
		StartPosition: 10,
		EndPosition:   10,
	}) {
		t.Errorf("Wrong token")
	}
}

func TestCopyVariableByLink(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const a = b;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	if !isSameToken(tokens[0], token.Token{
		Code:          token_variable_declaration.VARIABLE_DECLARAION,
		StartPosition: 0,
		EndPosition:   6,
	}) {
		t.Errorf("Wrong token")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_variable_declaration.VARIABLE_NAME_PARAM,
		Value:         "a",
		StartPosition: 6,
		EndPosition:   6,
	}) {
		t.Errorf("Wrong param")
	}

	if !isSameToken(tokens[1], token.Token{
		Code:          token.ASSIGNMENT,
		Value: "=",
		StartPosition: 8,
		EndPosition:   8,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[2], token.Token{
		Code:          token.KEY_WORD,
		Value:         "b",
		StartPosition: 10,
		EndPosition:   10,
	}) {
		t.Errorf("Wrong token")
	}
}

func TestVariableWithTokenName(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const functionA;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	if !isSameToken(tokens[0], token.Token{
		Code:          token_variable_declaration.VARIABLE_DECLARAION,
		StartPosition: 0,
		EndPosition:   14,
	}) {
		t.Errorf("Wrong token")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_variable_declaration.VARIABLE_NAME_PARAM,
		Value:         "functionA",
		StartPosition: 6,
		EndPosition:   14,
	}) {
		t.Errorf("Wrong param")
	}
}

func TestFunctionDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo() {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	if !isSameToken(tokens[0], token.Token{
		Code:          token_function_declaration.FUNCTION_DECLARATION,
		StartPosition: 0,
		EndPosition:   13,
	}) {
		t.Errorf("Wrong token")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_function_declaration.FUNCTION_NAME_PARAM,
		Value:         "foo",
		StartPosition: 9,
		EndPosition:   11,
	}) {
		t.Errorf("Wrong param")
	}

	if !isSameToken(tokens[1], token.Token{
		Code:          token.OPEN_BLOCK,
		Value: "{",
		StartPosition: 15,
		EndPosition:   15,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[2], token.Token{
		Code:          token.CLOSE_BLOCK,
		Value: "}",
		StartPosition: 16,
		EndPosition:   16,
	}) {
		t.Errorf("Wrong token")
	}
}

func TestFunctionDeclorationWithArguments(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo(a, b) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var tokens, _ = tokenizer.GetTokens()

	if !isSameToken(tokens[0], token.Token{
		Code:          token_function_declaration.FUNCTION_DECLARATION,
		StartPosition: 0,
		EndPosition:   17,
	}) {
		t.Errorf("Wrong token")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_function_declaration.FUNCTION_NAME_PARAM,
		Value:         "foo",
		StartPosition: 9,
		EndPosition:   11,
	}) {
		t.Errorf("Wrong param")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_function_declaration.FUNCTION_ARGUMENT_PARAM,
		Value:         "a",
		StartPosition: 13,
		EndPosition:   13,
	}) {
		t.Errorf("Wrong param")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_function_declaration.FUNCTION_ARGUMENT_PARAM,
		Value:         "b",
		StartPosition: 16,
		EndPosition:   16,
	}) {
		t.Errorf("Wrong param")
	}

	if !isSameToken(tokens[1], token.Token{
		Code:          token.OPEN_BLOCK,
		Value: "{",
		StartPosition: 19,
		EndPosition:   19,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[2], token.Token{
		Code:          token.CLOSE_BLOCK,
		Value: "}",
		StartPosition: 20,
		EndPosition:   20,
	}) {
		t.Errorf("Wrong token")
	}
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

	if !isSameToken(tokens[0], token.Token{
		Code:          token_function_declaration.FUNCTION_DECLARATION,
		StartPosition: 3,
		EndPosition:   16,
	}) {
		t.Errorf("Wrong token")
	}

	if !checkParam(tokens[0], token.TokenParam{
		Name:          token_function_declaration.FUNCTION_NAME_PARAM,
		Value:         "bar",
		StartPosition: 12,
		EndPosition:   14,
	}) {
		t.Errorf("Wrong param")
	}

	if !isSameToken(tokens[1], token.Token{
		Code:          token.OPEN_BLOCK,
		Value: "{",
		StartPosition: 18,
		EndPosition:   18,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[2], token.Token{
		Code:          token_return.RETURN_DECLARATION,
		StartPosition: 23,
		EndPosition:   28,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[3], token.Token{
		Code:          token_string.STRING,
		Value:         "Avada kedavra",
		StartPosition: 30,
		EndPosition:   44,
	}) {
		t.Errorf("Wrong token")
	}

	if !isSameToken(tokens[4], token.Token{
		Code:          token.CLOSE_BLOCK,
		Value: "}",
		StartPosition: 48,
		EndPosition:   48,
	}) {
		t.Errorf("Wrong token")
	}
}

func TestSyntaxError(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		const a = 123;

		^

		function foo() {

		}
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tokenizer = GetTokenizer(&buffer)
	var _, err = tokenizer.GetTokens()

	if err == nil {
		t.Errorf("Should return error")
	}

	re, ok := err.(parser_error.ParserError)

	if !ok {
		t.Errorf("Should return parser error")
	}

	if re.Message != "Syntax error, unexpected symbol: ^" {
		t.Errorf("Should return syntax error message, but receinve: \"%s\"", re.Message)
	}

	if re.StartPosition != 21 || re.EndPosition != 21 {
		t.Errorf("Should return position of error. Recived start: %d, end: %d", re.StartPosition, re.EndPosition)
	}
}

/// Utils

func isSameToken(first token.Token, second token.Token) bool {
	if first.Code != second.Code {
		fmt.Printf("Different token Codes: %s - %s\n", first.Code, second.Code)
		return false
	}

	if first.Value != second.Value {
		fmt.Printf("Different token Values: %s - %s\n", first.Value, second.Value)
		return false
	}

	if first.StartPosition != second.StartPosition {
		fmt.Printf("Different token start position: %d - %d\n", first.StartPosition, second.StartPosition)
		return false
	}

	if first.EndPosition != second.EndPosition {
		fmt.Printf("Different token end position: %d - %d\n", first.EndPosition, second.EndPosition)
		return false
	}

	return true
}

func checkParam(tokenItem token.Token, targetParam token.TokenParam) bool {
	for _, param := range tokenItem.Params {
		if param.Name == targetParam.Name && param.Value == targetParam.Value && param.StartPosition == targetParam.StartPosition && param.EndPosition == targetParam.EndPosition {
			return true
		}
	}

	fmt.Printf("Param not found in token")
	fmt.Printf("Looking for param - name: %s, value: %s, position, %d\n", targetParam.Name, targetParam.Value, targetParam.StartPosition)
	fmt.Printf("Token has params:")

	for _, param := range tokenItem.Params {
		fmt.Printf("name: %s, value: %s, position, %d\n", param.Name, param.Value, param.StartPosition)
	}

	return false
}
