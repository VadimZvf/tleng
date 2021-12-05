package token_keyword

import (
	"github.com/VadimZvf/golang/token"
)

var KEY_WORD = "KEY_WORD"
var KeyWordProcessor token.TokenProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if !token.IsKeyWordSymbol(buffer.GetSymbol()) {
		return token.Token{}, false, nil
	}

	var startPosition = buffer.GetPosition()

	for token.IsKeyWordSymbol(buffer.GetSymbol()) && !buffer.GetIsEnd() {
		buffer.AddSymbol()
		buffer.Next()
	}

	return token.Token{
		Code:          KEY_WORD,
		Value:         buffer.GetValue(),
		StartPosition: startPosition,
		EndPosition:   buffer.GetPosition() - 1,
		DebugValue:    buffer.GetValue(),
	}, true, nil
}
