package token_variable_decloration

import (
	"testing"

	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

type SourceMock struct {
	IsEnd    bool
	FullText string
	position int
}

func (source *SourceMock) NextSymbol() (symbol string, isEnd bool) {
	source.position += 1

	if source.position >= len(source.FullText) {
		return "", true
	}

	return string(source.FullText[source.position]), false
}

func GetSourceMock() *SourceMock {
	return &SourceMock{false, "", -1}
}

func TestReturnShouldNotBeFound(t *testing.T) {
	var src = GetSourceMock()
	src.FullText = `
		ret urn
	`

	var buffer = tokenizer_buffer.CreateBuffer(src)
	var token = token.Token{}
	var isFound = false

	for !buffer.GetIsEnd() && !isFound {
		token, isFound = ReturnProcessor(&buffer)
		buffer.TrimNext()
		buffer.Next()
		buffer.AddSymbol()
	}

	if isFound {
		t.Errorf("Should't find token")
	}

	if token.Code != "" {
		t.Errorf("Should't find token")
	}
}

func TestEmptyReturn(t *testing.T) {
	var src = GetSourceMock()
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
