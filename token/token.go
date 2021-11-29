package token

type TokenParam struct {
	Name     string
	Value    string
	Position int
}

type Token struct {
	Code     string
	Position int
	// Simple content like string content
	Value string
	// Detail information, like function arguments, return values etc...
	Params     []TokenParam
	DebugValue string
}

type IBuffer interface {
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

type TokenProcessor = func(buffer IBuffer) (foundToken Token, isFoundToken bool, err error)

func createSymbolProcessor(code string, symbol string) TokenProcessor {
	return func(buffer IBuffer) (foundToken Token, isFoundToken bool, err error) {
		if buffer.GetSymbol() == symbol {
			return Token{
				Code:       code,
				Position:   buffer.GetPosition(),
				DebugValue: buffer.GetSymbol(),
			}, true, nil
		}

		return Token{}, false, nil
	}
}

func createKeyWordProcessor(code string, keyWord string) TokenProcessor {
	return func(buffer IBuffer) (foundToken Token, isFoundToken bool, err error) {
		if buffer.GetValue() == keyWord {
			return Token{
				Code:       code,
				Position:   buffer.GetPosition(),
				DebugValue: buffer.GetValue(),
			}, true, nil
		}

		return Token{}, false, nil
	}
}

var NEW_LINE = "NEW_LINE"
var NewLineProcessor = createSymbolProcessor(NEW_LINE, "\n")

var ASSIGNMENT = "ASSIGNMENT"
var AssignmentProcessor = createSymbolProcessor(ASSIGNMENT, "=")

var OPEN_BLOCK = "OPEN_BLOCK"
var OpenBlockProcessor = createSymbolProcessor(OPEN_BLOCK, "{")

var CLOSE_BLOCK = "CLOSE_BLOCK"
var CloseBlockProcessor = createSymbolProcessor(CLOSE_BLOCK, "}")

var OPEN_EXPRESSION = "OPEN_EXPRESSION"
var OpenExpressionProcessor = createSymbolProcessor(OPEN_EXPRESSION, "(")

var CLOSE_EXPRESSION = "CLOSE_EXPRESSION"
var CloseExpressionProcessor = createSymbolProcessor(CLOSE_EXPRESSION, ")")

var END_LINE = "END_LINE"
var EndLineProcessor = createSymbolProcessor(END_LINE, ";")

var COMMA = "COMMA"
var CommaProcessor = createSymbolProcessor(COMMA, ",")

var DOT = "DOT"
var DotProcessor = createSymbolProcessor(DOT, ".")

var VARIABLE_DECLORAION = "VARIABLE_DECLORAION"
var VariableDeclorationProcessor = createKeyWordProcessor(VARIABLE_DECLORAION, "const")

var PROGRAMM = "PROGRAMM"
var KEY_WORD = "KEY_WORD"

// Utils
// =======================================
func IsLetter(s string) bool {
	for _, r := range s {
		if ((r < 'a' || r > 'z') && (r < 'A' || r > 'Z')) && r != '_' {
			return false
		}
	}
	return true
}

func ReadWord(buffer IBuffer) TokenParam {
	var position = buffer.GetPosition()
	for IsLetter(buffer.GetSymbol()) && !buffer.GetIsEnd() {
		position = buffer.GetPosition()
		buffer.AddSymbol()
		buffer.Next()
	}

	return TokenParam{
		Value:    buffer.GetValue(),
		Position: position,
	}
}
