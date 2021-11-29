package tokenizer_buffer

type iSource interface {
	NextSymbol() (symbol string, isEnd bool)
}

type Buffer struct {
	code     string
	value    string
	symbol   string
	source   iSource
	isEnd    bool
	position int
}

func (buffer *Buffer) GetValue() string {
	return buffer.value
}

func (buffer *Buffer) GetFullValue() string {
	return buffer.value + buffer.symbol
}

func (buffer *Buffer) GetSymbol() string {
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
	symbol, isEnd := buffer.source.NextSymbol()

	buffer.code = buffer.code + symbol
	buffer.symbol = symbol
	buffer.isEnd = isEnd
	buffer.position = buffer.position + 1
}

func (buffer *Buffer) TrimNext() {
	for (buffer.GetSymbol() == " " || buffer.GetSymbol() == "\t") && !buffer.GetIsEnd() {
		buffer.Next()
	}
}

func (buffer *Buffer) AddSymbol() {
	buffer.value = buffer.value + buffer.symbol
}

func (buffer *Buffer) Clear() {
	buffer.value = ""
}

func CreateBuffer(source iSource) Buffer {
	symbol, isEnd := source.NextSymbol()

	return Buffer{
		code:     symbol,
		value:    "",
		symbol:   symbol,
		source:   source,
		isEnd:    isEnd,
		position: 0,
	}
}
