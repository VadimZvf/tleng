package token_string

import (
	"fmt"

	"github.com/VadimZvf/golang/token"
)

var STRING = "STRING"
var StringProcessor token.TokenProcessor = proccess

func proccess(buffer token.IBuffer) (token.Token, bool) {
	if buffer.GetSymbol() != "\"" {
		return token.Token{}, false
	}

	fmt.Println("value ", buffer.GetValue())
	fmt.Println("symbol ", buffer.GetSymbol())

	// Remove quote mark at start
	buffer.Next()

	for buffer.GetSymbol() != "\"" && !buffer.GetIsEnd() {
		buffer.AddSymbol()
		buffer.Next()
	}

	stringToken := token.Token{
		Code:       STRING,
		Value:      buffer.GetValue(),
		Position:   buffer.GetPosition(),
		DebugValue: buffer.GetValue(),
	}

	// Remove quote mark at end
	buffer.Next()

	return stringToken, true
}
