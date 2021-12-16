package ast_token_stream

import "github.com/VadimZvf/golang/token"


type TokenStream struct {
	tokens []token.Token
	currentIndex int
}

func CreateTokenStream(tokens []token.Token) TokenStream {
	return TokenStream{
		tokens: tokens,
		currentIndex: 0,
	}
}

func (stream *TokenStream) Next() (token.Token, bool) {
	if stream.currentIndex + 1 >= len(stream.tokens) {
		return token.Token{}, true
	}

	var nextToken = stream.tokens[stream.currentIndex]

	stream.currentIndex = stream.currentIndex + 1

	return nextToken, false
}
