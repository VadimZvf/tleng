package tokenizer_buffer

type iSource interface {
	NextSymbol() (symbol rune, isEnd bool)
}

type Buffer struct {
	code     string
	value    string
	symbol   rune
	source   iSource
	isEnd    bool
	position int

	spiedSymbol rune
	isSpiedEnd  bool
}

func (buffer *Buffer) GetValue() string {
	return buffer.value
}

func (buffer *Buffer) GetFullValue() string {
	return buffer.value + string(buffer.symbol)
}

func (buffer *Buffer) GetSymbol() rune {
	return buffer.symbol
}

func (buffer *Buffer) GetPosition() int {
	return buffer.position
}

func (buffer *Buffer) GetIsEnd() bool {
	return buffer.isEnd
}

func (buffer *Buffer) GetReadedCode() string {
	return buffer.code
}

func (buffer *Buffer) Next() {
	buffer.position = buffer.position + 1

	if buffer.spiedSymbol != rune(0) {
		buffer.symbol = buffer.spiedSymbol
		buffer.spiedSymbol = rune(0)
		buffer.isEnd = buffer.isSpiedEnd
		return
	}

	symbol, isEnd := buffer.source.NextSymbol()

	buffer.code = buffer.code + string(symbol)
	buffer.symbol = symbol
	buffer.isEnd = isEnd

	if isEnd {
		buffer.symbol = rune(0)
	}
}

func (buffer *Buffer) TrimNext() {
	for (buffer.GetSymbol() == ' ' || buffer.GetSymbol() == '\t') && !buffer.GetIsEnd() {
		buffer.Next()
	}
}

func (buffer *Buffer) AddSymbol() {
	buffer.value = buffer.value + string(buffer.symbol)
}

func (buffer *Buffer) Clear() {
	buffer.value = ""
}

func (buffer *Buffer) PeekForward() rune {
	if buffer.spiedSymbol != rune(0) || buffer.isEnd || buffer.isSpiedEnd {
		return buffer.spiedSymbol
	}

	symbol, isEnd := buffer.source.NextSymbol()
	buffer.isSpiedEnd = isEnd
	buffer.spiedSymbol = symbol

	return buffer.spiedSymbol
}

func CreateBuffer(source iSource) Buffer {
	symbol, isEnd := source.NextSymbol()

	return Buffer{
		code:     string(symbol),
		value:    "",
		symbol:   symbol,
		source:   source,
		isEnd:    isEnd,
		position: 0,

		spiedSymbol: rune(0),
		isSpiedEnd:  isEnd,
	}
}
