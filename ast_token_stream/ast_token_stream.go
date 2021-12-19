package ast_token_stream

import "github.com/VadimZvf/golang/token"

type TokenStream struct {
	tokens       []token.Token
	currentIndex int
}

func CreateTokenStream(tokens []token.Token) TokenStream {
	return TokenStream{
		tokens:       tokens,
		currentIndex: 0,
	}
}

func (stream *TokenStream) MoveNext() {
	stream.currentIndex = stream.currentIndex + 1
}

func (stream *TokenStream) Look() (token.Token, bool) {
	if stream.currentIndex >= len(stream.tokens) {
		return token.Token{}, true
	}

	var token = stream.tokens[stream.currentIndex]

	return token, false
}

func (stream *TokenStream) LookNext() (token.Token, bool) {
	var nextIndex = stream.currentIndex + 1

	if nextIndex >= len(stream.tokens) {
		return token.Token{}, true
	}

	var nextToken = stream.tokens[nextIndex]

	return nextToken, false
}
