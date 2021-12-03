package token_variable_decloration

import (
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestVariableShouldNotBeFound(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `cons t = "";`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound, _ = VariableDeclarationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestEmptyVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const a;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound, _ = VariableDeclarationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if !isFound {
		t.Errorf("Should find token")
	}

	if foundToken.Code == "" {
		t.Errorf("Should find token")
	}

	var argument3Param = token.TokenParam{
		Name:     "NAME",
		Value:    "a",
		Position: 6,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save name")
	}
}

func TestLongNameVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const wow_foo_bar;`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound, _ = VariableDeclarationProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if !isFound {
		t.Errorf("Should find token")
	}

	if foundToken.Code == "" {
		t.Errorf("Should find token")
	}

	var argument3Param = token.TokenParam{
		Name:     "NAME",
		Value:    "wow_foo_bar",
		Position: 16,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save name")
	}
}

func TestErrorInvalidName(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `const 123foo`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = VariableDeclarationProcessor(&buffer)
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

	if re.Message != "Syntax error, variable cannot start with number" {
		t.Errorf("Should return variable name parsing error")
	}

	if re.Position != 0 || re.Length != 7 {
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
