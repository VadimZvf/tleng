package token_read_property

import (
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestReadProperty(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `foo.bar`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound, _ = ReadPropertyProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if token.Code != READ_PROPERTY {
		t.Errorf("Token should be found")
	}

	if token.StartPosition != 3 || token.EndPosition != 3 {
		t.Errorf("Should save token position")
	}
}

func TestReadPropertyWithUnderscore(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `foo._bar`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound, _ = ReadPropertyProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if token.Code != READ_PROPERTY {
		t.Errorf("Token should be found")
	}

	if token.StartPosition != 3 {
		t.Errorf("Should save token position")
	}
}

func TestErrorInvalidPropertyName(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `bar.%`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = ReadPropertyProcessor(&buffer)
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

	if re.Message != "Syntax error, invalid object property name" {
		t.Errorf("Should return syntax error. Recived: \"%s\"", re.Message)
	}

	if re.StartPosition != 3 || re.EndPosition != 4 {
		t.Errorf("Should return position of error. Recived start: %d, end: %d", re.StartPosition, re.EndPosition)
	}
}
