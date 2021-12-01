package token_function_declaration

import (
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestFunctionNotFound(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `func wow() {}
						wowdsa()
					`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound, _ = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound != false {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestFunctionWithoutArguments(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function wow() {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound, _ = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if foundToken.Code != "FUNCTION_DECLARATION" {
		t.Errorf("Token should be found")
	}

	if foundToken.Position != 7 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:     "NAME",
		Value:    "wow",
		Position: 11,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}
}

func TestIgnoreSpaces(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function      foo    (     ) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound, _ = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if foundToken.Code != "FUNCTION_DECLARATION" {
		t.Errorf("Token should be found")
	}

	if foundToken.Position != 7 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:     "NAME",
		Value:    "foo",
		Position: 16,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}
}

func TestParseSingleArgument(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function      bar(baz) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound, _ = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if foundToken.Code != "FUNCTION_DECLARATION" {
		t.Errorf("Token should be found")
	}

	if foundToken.Position != 7 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:     "NAME",
		Value:    "bar",
		Position: 16,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}

	var argumentParam = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "baz",
		Position: 20,
	}

	if !containParam(foundToken.Params, argumentParam) {
		t.Errorf("Should save argument")
	}
}

func TestParseManyArgument(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function      bar( baz,  foo,     gaz) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound, _ = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	var argumentParam = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "baz",
		Position: 21,
	}

	if !containParam(foundToken.Params, argumentParam) {
		t.Errorf("Should save argument")
	}

	var argument2Param = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "foo",
		Position: 27,
	}

	if !containParam(foundToken.Params, argument2Param) {
		t.Errorf("Should save argument")
	}

	var argument3Param = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "gaz",
		Position: 36,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save argument")
	}
}

func TestErrorNameParsing(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function ( baz,  foo,     gaz) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

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

	if re.Position != 0 || re.Length != 8 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.Position, re.Length)
	}
}

func TestErrorDeclorationParsing(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo baz,  foo,     gaz) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

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

	if re.Position != 0 || re.Length != 13 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.Position, re.Length)
	}
}

func TestErrorSecondBracketParsing(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo (gaz {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

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

	if re.Position != 0 || re.Length != 18 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.Position, re.Length)
	}
}

func TestErrorArgumentsParsing(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `function foo (  , a) {}`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

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

	if re.Position != 0 || re.Length != 16 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.Position, re.Length)
	}
}

func containParam(params []token.TokenParam, target token.TokenParam) bool {
	for _, param := range params {
		if param.Name == target.Name && param.Value == target.Value && param.Position == target.Position {
			return true
		}
	}
	return false
}
