package token

type TokenParam struct {
	Name          string
	Value         string
	StartPosition int
	EndPosition   int
}

type Token struct {
	Code          string
	StartPosition int
	EndPosition   int
	// Simple content like string content
	Value string
	// Detail information, like function arguments, return values etc...
	Params []TokenParam
}

type IBuffer interface {
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

type TokenProcessor = func(buffer IBuffer) (foundToken Token, isFoundToken bool, err error)

func createSymbolProcessor(code string, symbol rune) TokenProcessor {
	return func(buffer IBuffer) (foundToken Token, isFoundToken bool, err error) {
		if buffer.GetSymbol() == symbol {
			var position = buffer.GetPosition()
			buffer.Eat(1)
			return Token{
				Code:          code,
				StartPosition: position,
				EndPosition:   position,
				Value:         string(symbol),
			}, true, nil
		}

		return Token{}, false, nil
	}
}

var ASSIGNMENT = "ASSIGNMENT"
var AssignmentProcessor = createSymbolProcessor(ASSIGNMENT, '=')

var OPEN_BLOCK = "OPEN_BLOCK"
var OpenBlockProcessor = createSymbolProcessor(OPEN_BLOCK, '{')

var CLOSE_BLOCK = "CLOSE_BLOCK"
var CloseBlockProcessor = createSymbolProcessor(CLOSE_BLOCK, '}')

var OPEN_EXPRESSION = "OPEN_EXPRESSION"
var OpenExpressionProcessor = createSymbolProcessor(OPEN_EXPRESSION, '(')

var CLOSE_EXPRESSION = "CLOSE_EXPRESSION"
var CloseExpressionProcessor = createSymbolProcessor(CLOSE_EXPRESSION, ')')

var ADD = "ADD"
var AddProcessor = createSymbolProcessor(ADD, '+')

var END_LINE = "END_LINE"
var EndLineProcessor = createSymbolProcessor(END_LINE, ';')

var COMMA = "COMMA"
var CommaProcessor = createSymbolProcessor(COMMA, ',')

var PROGRAMM = "PROGRAMM"
var KEY_WORD = "KEY_WORD"

// Utils
// =======================================
func IsNumber(symbol rune) bool {
	return (symbol >= '0' && symbol <= '9')
}

func IsLetter(symbol rune) bool {
	return (symbol >= 'a' && symbol <= 'z') || (symbol >= 'A' && symbol <= 'Z')
}

func IsKeyWordSymbol(symbol rune) bool {
	return IsLetter(symbol) || IsNumber(symbol) || symbol == '_'
}

func IsValidKeyWord(s string) bool {
	for _, r := range s {
		if !IsKeyWordSymbol(r) {
			return false
		}
	}
	return true
}

func ReadWord(buffer IBuffer) TokenParam {
	var startPosition = buffer.GetPosition()

	for IsKeyWordSymbol(buffer.GetSymbol()) && !buffer.GetIsEnd() {
		buffer.AddSymbol()
		buffer.Next()
	}

	return TokenParam{
		Value:         buffer.GetValue(),
		StartPosition: startPosition,
		EndPosition:   buffer.GetPosition() - 1,
	}
}
