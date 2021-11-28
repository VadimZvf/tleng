package token_return

import (
	"github.com/VadimZvf/golang/token"
)

var RETURN_DECLARATION = "RETURN_DECLARATION"
var ReturnProcessor token.TokenProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool) {
	if buffer.GetFullValue() != "return" {
		return token.Token{}, false
	}

	var position = buffer.GetPosition()

	return token.Token{
		Code:     RETURN_DECLARATION,
		Position: position,
	}, true
}
