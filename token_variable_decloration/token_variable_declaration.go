package token_variable_decloration

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
)

var VARIABLE_DECLARAION = "VARIABLE_DECLARAION"
var VariableDeclarationProcessor token.TokenProcessor = proccess

var VARIABLE_NAME_PARAM = "NAME"

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if buffer.GetFullValue() != "const" || token.IsKeyWordSymbol(buffer.PeekForward()) {
		return token.Token{}, false, nil
	}

	var position = buffer.GetPosition()

	// Go to next symbol
	buffer.Next()
	buffer.TrimNext()
	buffer.Clear()

	var variableName = token.ReadWord(buffer)
	variableName.Name = VARIABLE_NAME_PARAM

	if len(variableName.Value) == 0 {
		// position-4, 5  Higlight text "const"
		return token.Token{}, false, parser_error.CreateError("Variable name should be defined!", position-4, 5)
	}

	return token.Token{
		Code:     VARIABLE_DECLARAION,
		Position: position,
		Params:   []token.TokenParam{variableName},
	}, true, nil
}
