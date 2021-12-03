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

	if len(buffer.GetValue()) == 0 {
		return token.Token{}, false, parser_error.ParserError{
			Message:  "Syntax error, unexpected dot",
			Position: buffer.GetPosition(),
			Length:   1,
		}
	}

	if !token.IsKeyWordSymbol(buffer.PeekForward()) {
		return token.Token{}, false, parser_error.ParserError{
			Message:  "Syntax error, invalid object property name",
			Position: buffer.GetPosition(),
			Length:   1,
		}
	}

	return token.Token{
		Code:       READ_PROPERTY,
		Position:   buffer.GetPosition(),
		DebugValue: ".",
	}, true, nil
}
