package tokenizer_buffer

import (
	"strings"
)

type iSource interface {
	NextSymbol() (symbol rune, isEnd bool)
}

type Buffer struct {
	value string // Current saved value, with trims, skips

	loadedValue      string // Current loaded value from source
	positionInBuffer int    // Index position in loaded value
	code             string // Full text of code
	position         int    // Current index position in code
	isEnd            bool

	// Inner props
	source iSource
}

func (buffer *Buffer) GetValue() string {
	return buffer.value
}

func (buffer *Buffer) GetSymbol() rune {
	if buffer.positionInBuffer >= len(buffer.loadedValue) {
		return rune(0)
	}

	return rune(buffer.loadedValue[buffer.positionInBuffer])
}

func (buffer *Buffer) GetPosition() int {
	return buffer.position + buffer.positionInBuffer
}

func (buffer *Buffer) GetIsEnd() bool {
	return buffer.positionInBuffer == len(buffer.loadedValue) && buffer.isEnd
}

func (buffer *Buffer) GetReadedCode() string {
	return buffer.code
}

func (buffer *Buffer) Next() {
	if buffer.GetIsEnd() {
		return
	}

	buffer.positionInBuffer = buffer.positionInBuffer + 1

	if buffer.positionInBuffer >= len(buffer.loadedValue) {
		buffer.loadSymbol()
	}
}

func (buffer *Buffer) Reset() {
	buffer.positionInBuffer = 0
}

func (buffer *Buffer) TrimNext() {
	for (buffer.GetSymbol() == ' ' || buffer.GetSymbol() == '\n' || buffer.GetSymbol() == '\t') && !buffer.GetIsEnd() {
		buffer.Next()
	}

	buffer.Eat(buffer.positionInBuffer)
}

func (buffer *Buffer) AddSymbol() {
	buffer.value = buffer.value + string(buffer.GetSymbol())
}

func (buffer *Buffer) IsStartsWith(tokenValue string) bool {
	for len(buffer.loadedValue) < len(tokenValue) && !buffer.isEnd {
		buffer.loadSymbol()
	}

	return strings.HasPrefix(buffer.loadedValue, tokenValue)
}

func (buffer *Buffer) Eat(length int) {
	buffer.loadedValue = buffer.loadedValue[length:]
	buffer.position = buffer.position + length
	buffer.positionInBuffer = 0
}

func (buffer *Buffer) Clear() {
	buffer.position = buffer.position + buffer.positionInBuffer
	buffer.loadedValue = buffer.loadedValue[buffer.positionInBuffer:]
	buffer.positionInBuffer = 0
	buffer.value = ""

	if len(buffer.loadedValue) == 0 {
		buffer.loadSymbol()
	}
}

func (buffer *Buffer) loadSymbol() {
	symbol, isEnd := buffer.source.NextSymbol()

	buffer.isEnd = isEnd

	if isEnd {
		return
	}

	buffer.code = buffer.code + string(symbol)
	buffer.loadedValue = buffer.loadedValue + string(symbol)
}

func CreateBuffer(source iSource) Buffer {
	symbol, isEnd := source.NextSymbol()

	return Buffer{
		value: "",

		loadedValue:      string(symbol),
		positionInBuffer: 0,
		code:             string(symbol),
		position:         0,
		isEnd:            isEnd,

		// Inner props
		source: source,
	}
}
