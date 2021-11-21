package token_function_declaration

import (
	"github.com/VadimZvf/golang/token"
)

var FUNCTION_DECLARATION = "FUNCTION_DECLARATION"
var FunctionDeclorationProcessor = proccess
var FUNCTION_NAME_PARAM = "NAME"
var FUNCTION_ARGUMENT_PARAM = "ARGUMENT"

func readWordsWithSeporator(buffer token.IBuffer, seporator string) []token.TokenParam {
	var words = []token.TokenParam{}
	buffer.TrimNext()

	for token.IsLetter(buffer.GetSymbol()) {
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

func proccess(buffer token.IBuffer) (token.Token, bool) {
	if buffer.GetValue() != "function" {
		return token.Token{}, false
	}

	var position = buffer.GetPosition()

	buffer.Clear()
	buffer.TrimNext()

	var functionName = token.ReadWord(buffer)
	functionName.Name = FUNCTION_NAME_PARAM

	buffer.Clear()
	buffer.TrimNext()

	if buffer.GetSymbol() == "(" {
		buffer.Next()
	}

	var arguments = readWordsWithSeporator(buffer, ",")

	for i := 0; i < len(arguments); i++ {
		param := &arguments[i]
		param.Name = FUNCTION_ARGUMENT_PARAM
	}

	if buffer.GetSymbol() == ")" {
		buffer.Next()
	}

	return token.Token{
		Code:     FUNCTION_DECLARATION,
		Position: position,
		Params:   append(arguments, functionName),
	}, true
}
