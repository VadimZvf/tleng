package token_function_declaration

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var FUNCTION_DECLARATION = "FUNCTION_DECLARATION"
var FunctionDeclorationProcessor = proccess
var FUNCTION_NAME_PARAM = "NAME"
var FUNCTION_ARGUMENT_PARAM = "ARGUMENT"

func readWordsWithSeporator(buffer token.IBuffer, seporator rune) []token.TokenParam {
	var words = []token.TokenParam{}
	buffer.TrimNext()

	for token.IsKeyWordSymbol(buffer.GetSymbol()) {
		words = append(words, token.ReadWord(buffer))
		buffer.Clear()
		buffer.TrimNext()
		if buffer.GetSymbol() != seporator {
			return words
		}
		buffer.Next()
		buffer.TrimNext()
	}

	return words
}

const functionDeclorationName = "function"

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if !buffer.IsStartsWith(functionDeclorationName + " ") {
		return token.Token{}, false, nil
	}

	var functionDeclorationStartPosition = buffer.GetPosition()

	buffer.Eat(len(functionDeclorationName))
	// Go to next symbol
	buffer.Next()

	buffer.Clear()
	buffer.TrimNext()

	var functionName = token.ReadWord(buffer)
	functionName.Name = FUNCTION_NAME_PARAM

	if len(functionName.Value) == 0 {
		return token.Token{}, false, parser_error.ParserError{
			Message:       "Function should have name",
			StartPosition: functionDeclorationStartPosition,
			EndPosition:   buffer.GetPosition(),
		}
	}

	buffer.Clear()
	buffer.TrimNext()

	if buffer.GetSymbol() != '(' {
		return token.Token{}, false, parser_error.ParserError{
			Message:       "Wrong function declaration syntax",
			StartPosition: functionDeclorationStartPosition,
			EndPosition:   buffer.GetPosition(),
		}
	}

	// Skip "("
	buffer.Next()

	var arguments = readWordsWithSeporator(buffer, ',')

	for i := 0; i < len(arguments); i++ {
		param := &arguments[i]
		param.Name = FUNCTION_ARGUMENT_PARAM

		if len(param.Value) == 0 {
			return token.Token{}, false, parser_error.ParserError{
				Message:       "Function argument should have name",
				StartPosition: functionDeclorationStartPosition,
				EndPosition:   buffer.GetPosition(),
			}
		}
	}

	if buffer.GetSymbol() != ')' {
		return token.Token{}, false, parser_error.ParserError{
			Message:       "Wrong function declaration syntax",
			StartPosition: functionDeclorationStartPosition,
			EndPosition:   buffer.GetPosition(),
		}
	}

	// Skip ")"
	buffer.Next()

	return token.Token{
		Code:          FUNCTION_DECLARATION,
		StartPosition: functionDeclorationStartPosition,
		EndPosition:   buffer.GetPosition() - 1,
		Params:        append(arguments, functionName),
	}, true, nil
}
