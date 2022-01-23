package token_string

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var STRING = "STRING"
var StringProcessor token.TokenProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if buffer.GetSymbol() != '"' && buffer.GetSymbol() != '\'' && buffer.GetSymbol() != '`' {
		return token.Token{}, false, nil
	}

	var stringWrapSymbol = buffer.GetSymbol()
	var startPosition = buffer.GetPosition()

	// Remove quote mark at start
	buffer.Next()

	for buffer.GetSymbol() != stringWrapSymbol {
		if stringWrapSymbol != '`' && buffer.GetSymbol() == '\n' {
			return token.Token{
					Code:          STRING,
					Value:         buffer.GetValue(),
					StartPosition: startPosition,
					EndPosition:   buffer.GetPosition(),
				}, false, parser_error.ParserError{
					Message:       "Syntax error, unexpected end of line. Use \"`\" for multiline string",
					StartPosition: startPosition,
					EndPosition:   buffer.GetPosition(),
				}
		}

		if buffer.GetIsEnd() {
			return token.Token{
					Code:          STRING,
					Value:         buffer.GetValue(),
					StartPosition: startPosition,
					EndPosition:   buffer.GetPosition(),
				}, false, parser_error.ParserError{
					Message:       "Syntax error, unexpected end of file",
					StartPosition: startPosition,
					EndPosition:   buffer.GetPosition(),
				}
		}

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
