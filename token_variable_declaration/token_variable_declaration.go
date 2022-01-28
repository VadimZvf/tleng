package token_variable_declaration

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var VARIABLE_DECLARAION = "VARIABLE_DECLARAION"
var VariableDeclarationProcessor token.TokenProcessor = proccess

var VARIABLE_NAME_PARAM = "NAME"
var variableDeclorationName = "var"

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if !buffer.IsStartsWithWord(variableDeclorationName) {
		return token.Token{}, false, nil
	}

	var startPosition = buffer.GetPosition()
	buffer.Eat(len(variableDeclorationName))

	buffer.Next()
	buffer.TrimNext()
	buffer.Clear()

	if token.IsNumber(buffer.GetSymbol()) {
		return token.Token{}, false, parser_error.ParserError{
			Message:       "Syntax error, variable cannot start with number",
			StartPosition: startPosition,
			EndPosition:   buffer.GetPosition(),
		}
	}

	var variableName = token.ReadWord(buffer)
	variableName.Name = VARIABLE_NAME_PARAM
	var endPosition = buffer.GetPosition() - 1

	if len(variableName.Value) == 0 {
		return token.Token{}, false, parser_error.ParserError{
			Message:       "Variable name should be defined!",
			StartPosition: startPosition,
			EndPosition:   endPosition,
		}
	}

	return token.Token{
		Code:          VARIABLE_DECLARAION,
		StartPosition: startPosition,
		EndPosition:   endPosition,
		Params:        []token.TokenParam{variableName},
	}, true, nil
}

func GetVariableNameParam(variableToken token.Token) token.TokenParam {
	for _, param := range variableToken.Params {
		if param.Name == VARIABLE_NAME_PARAM {
			return param
		}
	}

	return token.TokenParam{}
}
