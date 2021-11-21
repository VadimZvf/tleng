package token_function_declaration

import (
	"testing"

	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

type SourceMock struct {
	IsEnd    bool
	FullText string
	position int
}

func (source *SourceMock) NextSymbol() (symbol string, isEnd bool) {
	source.position += 1

	if source.position >= len(source.FullText) {
		return "", true
	}

	return string(source.FullText[source.position]), false
}

func GetSourceMock() *SourceMock {
	return &SourceMock{false, "", 0}
}

func TestFunctionNotFound(t *testing.T) {
	var src = GetSourceMock()
	src.FullText = `
		func wow() {

		}

		wowdsa()
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.Next()
		buffer.AddSymbol()
	}

	if isFound != false {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestFunctionWithoutArguments(t *testing.T) {
	var src = GetSourceMock()
	src.FullText = `
        function wow() {}
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound = FunctionDeclorationProcessor(&buffer)
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

	if foundToken.Position != 16 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:     "NAME",
		Value:    "wow",
		Position: 20,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}
}

func TestIgnoreSpaces(t *testing.T) {
	var src = GetSourceMock()
	src.FullText = `
        function      foo    (     ) {}
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound = FunctionDeclorationProcessor(&buffer)
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

	if foundToken.Position != 16 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:     "NAME",
		Value:    "foo",
		Position: 25,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}
}

func TestParseSingleArgument(t *testing.T) {
	var src = GetSourceMock()
	src.FullText = `
        function      bar(baz) {}
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound = FunctionDeclorationProcessor(&buffer)
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

	if foundToken.Position != 16 {
		t.Errorf("Should save token position")
	}

	var nameParam = token.TokenParam{
		Name:     "NAME",
		Value:    "bar",
		Position: 25,
	}

	if !containParam(foundToken.Params, nameParam) {
		t.Errorf("Should save function name")
	}

	var argumentParam = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "baz",
		Position: 29,
	}

	if !containParam(foundToken.Params, argumentParam) {
		t.Errorf("Should save argument")
	}
}

func TestParseManyArgument(t *testing.T) {
	var src = GetSourceMock()
	src.FullText = `
        function      bar( baz,  foo,     gaz) {}
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound = FunctionDeclorationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	var argumentParam = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "baz",
		Position: 30,
	}

	if !containParam(foundToken.Params, argumentParam) {
		t.Errorf("Should save argument")
	}

	var argument2Param = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "foo",
		Position: 36,
	}

	if !containParam(foundToken.Params, argument2Param) {
		t.Errorf("Should save argument")
	}

	var argument3Param = token.TokenParam{
		Name:     "ARGUMENT",
		Value:    "gaz",
		Position: 45,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save argument")
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
