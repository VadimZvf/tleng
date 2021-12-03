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

	for !buffer.GetIsEnd() && !isFound {
		token, isFound, _ = NumberProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

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

	for !buffer.GetIsEnd() && !isFound {
		token, isFound, _ = NumberProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if token.Code != NUMBER {
		t.Errorf("Token should be found")
	}

	if token.Position != 4 {
		t.Errorf("Should save token position")
	}
}

func TestErrorParsing(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `231wow`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var isFound = false
	var err = error(nil)

	for !buffer.GetIsEnd() && !isFound && err == nil {
		_, isFound, err = NumberProcessor(&buffer)
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

	if re.Message != "Syntax error, invalid keyword name. Keyword cannot start with number" {
		t.Errorf("Should return syntax error. Recived: \"%s\"", re.Message)
	}

	if re.Position != 0 || re.Length != 3 {
		t.Errorf("Should return position of error. Recived position: %d, length: %d", re.Position, re.Length)
	}
}
