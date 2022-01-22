package token_function_declaration

import (
	"fmt"
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestFunctionNotFound(t *testing.T) {
	var src = source_mock.GetSourceMock(`func wow() {}
	wowdsa()
`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := FunctionDeclorationProcessor(&buffer)

	if isFound != false {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestFunctionWithoutArguments(t *testing.T) {
	var src = source_mock.GetSourceMock(`function wow() {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	foundToken, isFound, _ := FunctionDeclorationProcessor(&buffer)

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if foundToken.Code != FUNCTION_DECLARATION {
		t.Errorf("Token should be found")
	}

	if foundToken.StartPosition != 0 || foundToken.EndPosition != 13 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:          FUNCTION_NAME_PARAM,
		Value:         "wow",
		StartPosition: 9,
		EndPosition:   11,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}
}

func TestIgnoreSpaces(t *testing.T) {
	var src = source_mock.GetSourceMock(`function      foo    (     ) {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	foundToken, isFound, _ := FunctionDeclorationProcessor(&buffer)

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if foundToken.Code != FUNCTION_DECLARATION {
		t.Errorf("Token should be found")
	}

	if foundToken.StartPosition != 0 || foundToken.EndPosition != 27 {
		t.Errorf("Should save token position. Received start: %d end: %d", foundToken.StartPosition, foundToken.EndPosition)
	}

	var nameParam = token.TokenParam{
		Name:          FUNCTION_NAME_PARAM,
		Value:         "foo",
		StartPosition: 14,
		EndPosition:   16,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}
}

func TestParseSingleArgument(t *testing.T) {
	var src = source_mock.GetSourceMock(`function      bar(baz) {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	foundToken, isFound, _ := FunctionDeclorationProcessor(&buffer)

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if foundToken.Code != FUNCTION_DECLARATION {
		t.Errorf("Token should be found")
	}

	if foundToken.StartPosition != 0 || foundToken.EndPosition != 21 {
		t.Errorf("Should save token position. Received start: %d end: %d", foundToken.StartPosition, foundToken.EndPosition)
	}

	var nameParam = token.TokenParam{
		Name:          FUNCTION_NAME_PARAM,
		Value:         "bar",
		StartPosition: 14,
		EndPosition:   16,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}

	var argumentParam = token.TokenParam{
		Name:          FUNCTION_ARGUMENT_PARAM,
		Value:         "baz",
		StartPosition: 18,
		EndPosition:   20,
	}

	if !containParam(foundToken.Params, argumentParam) {
		t.Errorf("Should save argument")
	}
}

func TestParseManyArgument(t *testing.T) {
	var src = source_mock.GetSourceMock(`function      bar( baz,  foo,     gaz) {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	foundToken, _, _ := FunctionDeclorationProcessor(&buffer)

	var argumentParam = token.TokenParam{
		Name:          FUNCTION_ARGUMENT_PARAM,
		Value:         "baz",
		StartPosition: 19,
		EndPosition:   21,
	}

	if !containParam(foundToken.Params, argumentParam) {
		t.Errorf("Should save argument")
	}

	var argument2Param = token.TokenParam{
		Name:          FUNCTION_ARGUMENT_PARAM,
		Value:         "foo",
		StartPosition: 25,
		EndPosition:   27,
	}

	if !containParam(foundToken.Params, argument2Param) {
		t.Errorf("Should save argument")
	}

	var argument3Param = token.TokenParam{
		Name:          FUNCTION_ARGUMENT_PARAM,
		Value:         "gaz",
		StartPosition: 34,
		EndPosition:   36,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save argument")
	}
}

func TestErrorNameParsing(t *testing.T) {
	var src = source_mock.GetSourceMock(`function ( baz,  foo,     gaz) {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var err = error(nil)

	_, _, err = FunctionDeclorationProcessor(&buffer)

	if err == nil {
		t.Errorf("Should return error")
	}

	re, ok := err.(parser_error.ParserError)

	if !ok {
		t.Errorf("Should return parser error")
	}

	if re.Message != "Function should have name" {
		t.Errorf("Should return function name parsing error")
	}

	if re.StartPosition != 0 || re.EndPosition != 9 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.StartPosition, re.EndPosition)
	}
}

func TestErrorDeclorationParsing(t *testing.T) {
	var src = source_mock.GetSourceMock(`function foo baz,  foo,     gaz) {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var err = error(nil)

	_, _, err = FunctionDeclorationProcessor(&buffer)

	if err == nil {
		t.Errorf("Should return error")
	}

	re, ok := err.(parser_error.ParserError)

	if !ok {
		t.Errorf("Should return parser error")
	}

	if re.Message != "Wrong function declaration syntax" {
		t.Errorf("Should return function parsing error. Recived: \"%s\"", re.Message)
	}

	if re.StartPosition != 0 || re.EndPosition != 13 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.StartPosition, re.EndPosition)
	}
}

func TestErrorSecondBracketParsing(t *testing.T) {
	var src = source_mock.GetSourceMock(`function foo (gaz {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var err = error(nil)

	_, _, err = FunctionDeclorationProcessor(&buffer)

	if err == nil {
		t.Errorf("Should return error")
	}

	re, ok := err.(parser_error.ParserError)

	if !ok {
		t.Errorf("Should return parser error")
	}

	if re.Message != "Wrong function declaration syntax" {
		t.Errorf("Should return function parsing error. Recived: \"%s\"", re.Message)
	}

	if re.StartPosition != 0 || re.EndPosition != 18 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.StartPosition, re.EndPosition)
	}
}

func TestErrorArgumentsParsing(t *testing.T) {
	var src = source_mock.GetSourceMock(`function foo (  , a) {}`)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var err = error(nil)

	_, _, err = FunctionDeclorationProcessor(&buffer)

	if err == nil {
		t.Errorf("Should return error")
	}

	re, ok := err.(parser_error.ParserError)

	if !ok {
		t.Errorf("Should return parser error")
	}

	if re.Message != "Wrong function declaration syntax" {
		t.Errorf("Should return function parsing error. Recived: \"%s\"", re.Message)
	}

	if re.StartPosition != 0 || re.EndPosition != 16 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.StartPosition, re.EndPosition)
	}
}

func containParam(params []token.TokenParam, target token.TokenParam) bool {
	for _, param := range params {
		if isSameParams(param, target) {
			return true
		}
	}

	fmt.Println("Params:", params)
	fmt.Println("Target:", target)
	return false
}

func isSameParams(param token.TokenParam, target token.TokenParam) bool {
	if param.Name != target.Name {
		return false
	}

	if param.Value != target.Value {
		return false
	}

	if param.StartPosition != target.StartPosition {
		return false
	}

	if param.EndPosition != target.EndPosition {
		return false
	}

	return true
}
