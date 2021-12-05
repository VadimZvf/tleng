package token_read_property

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var READ_PROPERTY = "READ_PROPERTY"
var ReadPropertyProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if buffer.GetSymbol() != '.' {
		return token.Token{}, false, nil
	}

	var startPosition = buffer.GetPosition()
	buffer.Next()

	if !token.IsKeyWordSymbol(buffer.GetSymbol()) {
		return token.Token{}, false, parser_error.ParserError{
			Message:       "Syntax error, invalid object property name",
			StartPosition: startPosition,
			EndPosition:   buffer.GetPosition(),
		}
	}

	return token.Token{
		Code:          READ_PROPERTY,
		StartPosition: startPosition,
		EndPosition:   startPosition,
		DebugValue:    ".",
	}, true, nil
}
