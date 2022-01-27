package token_boolean

import (
	"github.com/VadimZvf/golang/token"
)

var BOOLEAN = "BOOLEAN"
var BooleanProcessor token.TokenProcessor = proccess

var booleanTrueValue = "true"
var booleanFalseValue = "false"

func proccess(buffer token.IBuffer) (token.Token, bool, error) {
	if buffer.IsStartsWith(booleanTrueValue) {
		return proccessTrue(buffer)
	}

	if buffer.IsStartsWith(booleanFalseValue) {
		return proccessFalse(buffer)
	}

	return token.Token{}, false, nil
}

func proccessTrue(buffer token.IBuffer) (token.Token, bool, error) {
	var startPosition = buffer.GetPosition()

	buffer.Eat(len(booleanTrueValue))

	var endPosition = buffer.GetPosition()

	buffer.Clear()

	return token.Token{
		Code:          BOOLEAN,
		Value:         booleanTrueValue,
		StartPosition: startPosition,
		EndPosition:   endPosition - 1,
	}, true, nil
}

func proccessFalse(buffer token.IBuffer) (token.Token, bool, error) {
	var startPosition = buffer.GetPosition()

	buffer.Eat(len(booleanFalseValue))

	var endPosition = buffer.GetPosition()

	buffer.Clear()

	return token.Token{
		Code:          BOOLEAN,
		Value:         booleanFalseValue,
		StartPosition: startPosition,
		EndPosition:   endPosition - 1,
	}, true, nil
}
