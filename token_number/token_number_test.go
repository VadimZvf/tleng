package token_number

import (
	"testing"

	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestNumberNotFound(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `foo323`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	token, isFound, _ = NumberProcessor(&buffer)

	if isFound != false {
		t.Errorf("Token should'nt be found")
	}

	if token.Code == NUMBER {
		t.Errorf("Token should'nt be found")
	}
}

func TestNumber(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `32112`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	token, isFound, _ = NumberProcessor(&buffer)

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if token.Code != NUMBER {
		t.Errorf("Token should be found")
	}

	if token.StartPosition != 0 || token.EndPosition != 4 {
		t.Errorf("Should save token position")
	}
}

func TestErrorParsing(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `231wow`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var err = error(nil)

	_, _, err = NumberProcessor(&buffer)

	if err == nil {
		t.Errorf("Should return error")
	}

	re, ok := err.(parser_error.ParserError)

	if !ok {
		t.Errorf("Should return parser error")
	}

	if re.Message != "Syntax error, invalid keyword name. Keyword cannot start with number" {
		t.Errorf("Should return syntax error. Recived: \"%s\"", re.Message)
	}

	if re.StartPosition != 0 || re.EndPosition != 3 {
		t.Errorf("Should return position of error. Recived start: %d, end: %d", re.StartPosition, re.EndPosition)
	}
}
