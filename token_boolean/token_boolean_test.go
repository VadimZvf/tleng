package token_boolean

import (
	"testing"

	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestBooleanNotFound(t *testing.T) {
	var src = source_mock.GetSourceMock(`truefalseVariable`)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	token, isFound, _ = BooleanProcessor(&buffer)

	if isFound != false {
		t.Errorf("Token should'nt be found")
	}

	if token.Code == BOOLEAN {
		t.Errorf("Token should'nt be found")
	}
}

func TestTrue(t *testing.T) {
	var src = source_mock.GetSourceMock(`true `)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	token, isFound, _ = BooleanProcessor(&buffer)

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if token.Code != BOOLEAN {
		t.Errorf("Token should be found")
	}

	if token.StartPosition != 0 || token.EndPosition != 3 {
		t.Errorf("Should save token position")
	}
}

func TestFalse(t *testing.T) {
	var src = source_mock.GetSourceMock(`false`)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	token, isFound, _ = BooleanProcessor(&buffer)

	if isFound == false {
		t.Errorf("Token should be found")
	}

	if token.Code != BOOLEAN {
		t.Errorf("Token should be found")
	}

	if token.Value != "false" {
		t.Errorf("Should save token value")
	}

	if token.StartPosition != 0 || token.EndPosition != 4 {
		t.Errorf("Should save token position")
	}
}
