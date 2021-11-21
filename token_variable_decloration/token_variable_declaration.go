package token_variable_decloration

import (
	"github.com/VadimZvf/golang/token"
)

var VARIABLE_DECLARAION = "VARIABLE_DECLARAION"
var VariableDeclarationProcessor token.TokenProcessor = proccess

var VARIABLE_NAME_PARAM = "NAME"

func proccess(buffer token.IBuffer) (token.Token, bool) {
	if buffer.GetValue() != "const" {
		return token.Token{}, false
	}

	var position = buffer.GetPosition()

	buffer.TrimNext()
	buffer.Clear()

	var variableName = token.ReadWord(buffer)
	variableName.Name = VARIABLE_NAME_PARAM

	return token.Token{
		Code:     VARIABLE_DECLARAION,
		Position: position,
		Params:   []token.TokenParam{variableName},
	}, true
}
