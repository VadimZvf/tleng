package token_variable_declaration

import (
	"fmt"
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestVariableShouldNotBeFound(t *testing.T) {
	var src = source_mock.GetSourceMock(`va r = "";`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := VariableDeclarationProcessor(&buffer)

	if isFound {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestEmptyVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock(`var a;`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	foundToken, isFound, _ := VariableDeclarationProcessor(&buffer)

	if !isFound {
		t.Errorf("Should find token")
	}

	if foundToken.Code == "" {
		t.Errorf("Should find token")
	}

	var variableNameParam = token.TokenParam{
		Name:          "NAME",
		Value:         "a",
		StartPosition: 4,
		EndPosition:   4,
	}

	if !containParam(foundToken.Params, variableNameParam) {
		t.Errorf("Should save name")
	}
}

func TestLongNameVariableDecloration(t *testing.T) {
	var src = source_mock.GetSourceMock(`var wow_foo_bar;`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	foundToken, isFound, _ := VariableDeclarationProcessor(&buffer)

	if !isFound {
		t.Errorf("Should find token")
	}

	if foundToken.Code == "" {
		t.Errorf("Should find token")
	}

	var argument3Param = token.TokenParam{
		Name:          "NAME",
		Value:         "wow_foo_bar",
		StartPosition: 4,
		EndPosition:   14,
	}

	if !containParam(foundToken.Params, argument3Param) {
		t.Errorf("Should save name")
	}
}

func TestErrorInvalidName(t *testing.T) {
	var src = source_mock.GetSourceMock(`var 123foo`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	_, _, err := VariableDeclarationProcessor(&buffer)

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

	if re.StartPosition != 0 || re.EndPosition != 4 {
		t.Errorf("Should return position of error. Recived start: %d, end: %d", re.StartPosition, re.EndPosition)
	}
}

func containParam(params []token.TokenParam, target token.TokenParam) bool {
	for _, param := range params {
		if param.Name != target.Name {
			fmt.Printf("Invalid property name. Expected: %s Received: %s\n", target.Name, param.Name)
			return false
		}

		if param.Value != target.Value {
			fmt.Printf("Invalid property value. Expected: %s Received: %s\n", target.Value, param.Value)
			return false
		}

		if param.StartPosition != target.StartPosition {
			fmt.Printf("Invalid property StartPosition. Expected: %d Received: %d\n", target.StartPosition, param.StartPosition)
			return false
		}

		if param.EndPosition != target.EndPosition {
			fmt.Printf("Invalid property EndPosition. Expected: %d Received: %d\n", target.EndPosition, param.EndPosition)
			return false
		}
	}
	return true
}
