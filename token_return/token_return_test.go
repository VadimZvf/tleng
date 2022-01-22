package token_return

import (
	"testing"

	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestReturnShouldNotBeFound(t *testing.T) {
	var src = source_mock.GetSourceMock(`returnfoo`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := ReturnProcessor(&buffer)

	if isFound {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestEmptyReturn(t *testing.T) {
	var src = source_mock.GetSourceMock(`return`)
	var buffer = tokenizer_buffer.CreateBuffer(src)

	token, isFound, _ := ReturnProcessor(&buffer)

	if !isFound {
		t.Errorf("Should find token")
	}

	if token.Code != RETURN_DECLARATION {
		t.Errorf("Should find token")
	}

	if token.StartPosition != 0 || token.EndPosition != 5 {
		t.Errorf("Should save position. But receive start: %d end: %d", token.StartPosition, token.EndPosition)
	}
}
