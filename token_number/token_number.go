package token_number

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var NUMBER = "NUMBER"
var NumberProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if !token.IsNumber(buffer.GetSymbol()) {
		return token.Token{}, false, nil
	}

	if len(buffer.GetValue()) != 0 {
		return token.Token{}, false, nil
	}

	var startPosition = buffer.GetPosition()

	for token.IsNumber(buffer.GetSymbol()) && !buffer.GetIsEnd() {
		buffer.AddSymbol()
		buffer.Next()
	}

	if token.IsKeyWordSymbol(buffer.GetSymbol()) {
		return token.Token{}, false, parser_error.ParserError{
			Message:  "Syntax error, invalid keyword name. Keyword cannot start with number",
			Position: startPosition,
			Length:   buffer.GetPosition() - startPosition,
		}
	}

	return token.Token{
		Code:       NUMBER,
		Position:   buffer.GetPosition() - 1,
		Value:      buffer.GetValue(),
		DebugValue: buffer.GetValue(),
	}, true, nil
}
