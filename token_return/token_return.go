package token_return

import (
	"github.com/VadimZvf/golang/token"
)

var RETURN_DECLARATION = "RETURN_DECLARATION"
var ReturnProcessor token.TokenProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if buffer.GetFullValue() != "return" || token.IsKeyWordSymbol(buffer.PeekForward()) {
		return token.Token{}, false, nil
	}

	var position = buffer.GetPosition()

	return token.Token{
		Code:     RETURN_DECLARATION,
		Position: position,
	}, true, nil
}
