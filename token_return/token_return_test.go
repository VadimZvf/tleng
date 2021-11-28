package token_return

import (
	"testing"

	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

func TestReturnShouldNotBeFound(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
		ret urn
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound = ReturnProcessor(&buffer)
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

func TestEmptyReturn(t *testing.T) {
	var src = source_mock.GetSourceMock()
	src.FullText = `
        return
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound = ReturnProcessor(&buffer)
		buffer.TrimNext()
		buffer.AddSymbol()
		buffer.Next()
	}

	if !isFound {
		t.Errorf("Should find token")
	}

	if token.Code != "RETURN_DECLARATION" {
		t.Errorf("Should find token")
	}

	if token.Position != 15 {
		t.Errorf("Should save position")
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
