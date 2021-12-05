package token_return

import (
	"github.com/VadimZvf/golang/token"
)

var RETURN_DECLARATION = "RETURN_DECLARATION"
var ReturnProcessor token.TokenProcessor = proccess
var returnName = "return"

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if !buffer.IsStartsWith(returnName) {
		return token.Token{}, false, nil
	}

	var startPosition = buffer.GetPosition()

	for i := 0; i < len(returnName); i++ {
		buffer.Next()
	}

	if token.IsKeyWordSymbol(buffer.GetSymbol()) {
		return token.Token{}, false, nil
	}

	buffer.Eat(len(returnName))

	return token.Token{
		Code:          RETURN_DECLARATION,
		StartPosition: startPosition,
		EndPosition:   buffer.GetPosition() - 1,
	}, true, nil
}
