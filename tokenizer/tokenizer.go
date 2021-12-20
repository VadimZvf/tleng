package tokenizer

import (
	"github.com/VadimZvf/golang/parser_error"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_keyword"
	"github.com/VadimZvf/golang/token_number"
	"github.com/VadimZvf/golang/token_read_property"
	"github.com/VadimZvf/golang/token_return"
	"github.com/VadimZvf/golang/token_string"
	"github.com/VadimZvf/golang/token_variable_declaration"
)

type iBuffer interface {
	GetValue() (value string)
	GetSymbol() (symbol rune)
	GetPosition() int
	GetIsEnd() bool
	Next()
	TrimNext()
	AddSymbol()
	IsStartsWith(value string) bool
	Eat(length int)
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

		foundToken, isFoundToken, err := getToken(tknzr.buffer)

		if err != nil {
			return tknzr.tokens, err
		}

		if isFoundToken {
			tknzr.addToken(foundToken)
			tknzr.buffer.Clear()
		}

		if tknzr.buffer.GetIsEnd() {
			return tknzr.tokens, nil
		}

		if !isFoundToken {
			return tknzr.tokens, parser_error.ParserError{
				Message:       "Syntax error, unexpected symbol: " + string(tknzr.buffer.GetSymbol()),
				StartPosition: tknzr.buffer.GetPosition(),
				EndPosition:   tknzr.buffer.GetPosition(),
			}
		}
	}
}

func getToken(buffer iBuffer) (token.Token, bool, error) {
	var tokensArray = []token.TokenProcessor{
		token_read_property.ReadPropertyProcessor,
		token_number.NumberProcessor,
		token_return.ReturnProcessor,
		token_variable_declaration.VariableDeclarationProcessor,
		token_function_declaration.FunctionDeclorationProcessor,
		token_keyword.KeyWordProcessor,
		token_string.StringProcessor,
		token.AssignmentProcessor,
		token.OpenBlockProcessor,
		token.CloseBlockProcessor,
		token.OpenExpressionProcessor,
		token.CloseExpressionProcessor,
		token.AddProcessor,
		token.SubtractProcessor,
		token.EndLineProcessor,
		token.CommaProcessor,
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
