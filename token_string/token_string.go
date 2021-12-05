package token_string

import (
	"github.com/VadimZvf/golang/token"
)

var STRING = "STRING"
var StringProcessor token.TokenProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if buffer.GetSymbol() != '"' {
		return token.Token{}, false, nil
	}

	var startPosition = buffer.GetPosition()

	// Remove quote mark at start
	buffer.Next()

	for buffer.GetSymbol() != '"' && !buffer.GetIsEnd() {
		buffer.AddSymbol()
		buffer.Next()
	}

	stringToken := token.Token{
		Code:          STRING,
		Value:         buffer.GetValue(),
		StartPosition: startPosition,
		EndPosition:   buffer.GetPosition(),
	}

	// Remove quote mark at end
	buffer.Next()

	return stringToken, true, nil
}
