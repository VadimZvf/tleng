package token_string

import (
	"testing"

	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestDoubleQuote(t *testing.T) {
	var src = source_mock.GetSourceMock(`"some text"`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := StringProcessor(&buffer)

	if !isFound {
		t.Errorf("Should find token")
	}

	if token.Code != STRING {
		t.Errorf("Should find token")
	}

	if token.Value != "some text" {
		t.Errorf("Should save value")
	}
}

func TestSingleQuote(t *testing.T) {
	var src = source_mock.GetSourceMock(`'some text with single quote'`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := StringProcessor(&buffer)

	if !isFound {
		t.Errorf("Should find token")
	}

	if token.Code != STRING {
		t.Errorf("Should find token")
	}

	if token.Value != "some text with single quote" {
		t.Errorf("Should save value")
	}
}

func TestBackQuote(t *testing.T) {
	var valueString = `some text
	with back\n
	quote`
	var src = source_mock.GetSourceMock("`" + valueString + "`")
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := StringProcessor(&buffer)

	if !isFound {
		t.Errorf("Should find token")
	}

	if token.Code != STRING {
		t.Errorf("Should find token")
	}

	if token.Value != valueString {
		t.Errorf("Should save value")
	}
}
