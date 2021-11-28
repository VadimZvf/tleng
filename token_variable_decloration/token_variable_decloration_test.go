package token_variable_decloration

import (
	"testing"

	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestVariableShouldNotBeFound(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		cons t = "";
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound = VariableDeclarationProcessor(&buffer)
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
	src.FullText = `
		const a;
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound = VariableDeclarationProcessor(&buffer)
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
		Position: 10,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save name")
	}
}

func TestLongNameVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		const wow_foo_bar;
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var foundToken = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		foundToken, isFound = VariableDeclarationProcessor(&buffer)
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
		Position: 20,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save name")
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
