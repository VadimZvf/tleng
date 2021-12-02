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
	GetSymbol() (symbol rune)
	GetPosition() int
	GetIsEnd() bool
	Next()
	TrimNext()
	AddSymbol()
	PeekForward() rune
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

func (tknzr *Tokenizer) GetTokens() ([]token.Token, error) {
	for {
		tknzr.buffer.TrimNext()

		keyWordToken, isFoundKeyWordToken, err := getKeyWordToken(tknzr.buffer)

		if err != nil {
			return tknzr.tokens, err
		}

		if isFoundKeyWordToken {
			tknzr.addToken(keyWordToken)
			tknzr.buffer.Clear()
		}

		symbolToken, isFoundSymbolToken := getSymbolToken(tknzr.buffer)

		if !isFoundKeyWordToken && !isFoundSymbolToken {
			tknzr.buffer.AddSymbol()
		}

		// Unknown token, maybe its reference to variable
		if (tknzr.buffer.GetIsEnd() || isFoundSymbolToken) && len(tknzr.buffer.GetValue()) > 0 {
			tknzr.addToken(getUnknownKeyWordToken(tknzr.buffer))
		}

		if isFoundSymbolToken {
			tknzr.addToken(symbolToken)
			tknzr.buffer.Clear()
		}

		if tknzr.buffer.GetIsEnd() {
			return tknzr.tokens, nil
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
		token, isFound, _ := tokenProccessor(buffer)

		if isFound {
			return token, true
		}
	}

	return token.Token{}, false
}

func getKeyWordToken(buffer iBuffer) (token.Token, bool, error) {
	var tokensArray = []token.TokenProcessor{
		token_return.ReturnProcessor,
		token_variable_decloration.VariableDeclarationProcessor,
		token_function_declaration.FunctionDeclorationProcessor,
		token_string.StringProcessor,
	}

	for i := 0; i < len(tokensArray); i++ {
		tokenProccessor := tokensArray[i]
		token, isFound, err := tokenProccessor(buffer)

		if err != nil {
			return token, false, err
		}

		if isFound {
			return token, true, nil
		}
	}

	return token.Token{}, false, nil
}

func getUnknownKeyWordToken(buffer iBuffer) token.Token {
	return token.Token{
		Code:       token.KEY_WORD,
		Value:      buffer.GetValue(),
		DebugValue: buffer.GetValue(),
		// Substract 1 becouse keyword may be found only next token
		// thats mean current position at space or END_LINK token
		Position: buffer.GetPosition() - 1,
	}
}
