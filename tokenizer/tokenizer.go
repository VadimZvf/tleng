package tokenizer

import (
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
	"github.com/VadimZvf/golang/token_variable_decloration"
)

type iBuffer interface {
	GetValue() (value string)
	GetFullValue() (value string)
	GetSymbol() (symbol string)
	GetPosition() int
	GetIsEnd() bool
	Next()
	TrimNext()
	AddSymbol()
	Clear()
}

type Tokenizer struct {
	buffer iBuffer
	tokens []token.Token
	token  token.Token
}

func GetTokenizer(buffer iBuffer) Tokenizer {
	return Tokenizer{
		buffer: buffer,
		tokens: []token.Token{},
		token:  token.Token{},
	}
}

func (tknzr *Tokenizer) addToken(token token.Token) {
	tknzr.tokens = append(tknzr.tokens, token)
}

func (tknzr *Tokenizer) GetTokens() []token.Token {
	for {
		tknzr.buffer.TrimNext()

		keyWordToken, isFoundKeyWordToken := getKeyWordToken(tknzr.buffer)

		if isFoundKeyWordToken {
			tknzr.addToken(keyWordToken)
			tknzr.buffer.Clear()
		}

		symbolToken, isFoundSymbolToken := getSymbolToken(tknzr.buffer)

		if isFoundSymbolToken {
			tknzr.addToken(symbolToken)
			tknzr.buffer.Clear()
		}

		if !isFoundKeyWordToken && !isFoundSymbolToken {
			tknzr.buffer.AddSymbol()
		}

		// Unknown token, maybe its reference to variable
		if (tknzr.buffer.GetIsEnd() || isFoundSymbolToken) && len(tknzr.buffer.GetValue()) > 0 {
			tknzr.addToken(getUnknownKeyWordToken(tknzr.buffer))
		}

		if tknzr.buffer.GetIsEnd() {
			return tknzr.tokens
		}

		tknzr.buffer.Next()
	}
}

func getSymbolToken(buffer iBuffer) (token.Token, bool) {
	var tokensArray = []token.TokenProcessor{
		token.NewLineProcessor,
		token.AssignmentProcessor,
		token.OpenBlockProcessor,
		token.CloseBlockProcessor,
		token.OpenExpressionProcessor,
		token.CloseExpressionProcessor,
		token.EndLineProcessor,
		token.DotProcessor,
		token.CommaProcessor,
	}

	for i := 0; i < len(tokensArray); i++ {
		tokenProccessor := tokensArray[i]
		token, isFound := tokenProccessor(buffer)

		if isFound {
			return token, true
		}
	}

	return token.Token{}, false
}

func getKeyWordToken(buffer iBuffer) (token.Token, bool) {
	var tokensArray = []token.TokenProcessor{
		token_return.ReturnProcessor,
		token_variable_decloration.VariableDeclarationProcessor,
		token_function_declaration.FunctionDeclorationProcessor,
		token_string.StringProcessor,
	}

	for i := 0; i < len(tokensArray); i++ {
		tokenProccessor := tokensArray[i]
		token, isFound := tokenProccessor(buffer)

		if isFound {
			return token, true
		}
	}

	return token.Token{}, false
}

func getUnknownKeyWordToken(buffer iBuffer) token.Token {
	return token.Token{
		Code:       token.KEY_WORD,
		Value:      buffer.GetValue(),
		DebugValue: buffer.GetValue(),
		Position:   buffer.GetPosition(),
	}
}
