package source_mock

type SimpleSourceMock struct {
	IsEnd           bool
	NextSymbolValue rune
}

func (source *SimpleSourceMock) NextSymbol() (rune, bool) {
	return source.NextSymbolValue, source.IsEnd
}

func GetSimpleSource() *SimpleSourceMock {
	return &SimpleSourceMock{false, rune(0)}
}

type SourceMock struct {
	IsEnd    bool
	FullText string
	position int
}

func (source *SourceMock) NextSymbol() (rune, bool) {
	source.position += 1

	if source.position >= len(source.FullText) {
		return rune(0), true
	}

	var symbol = source.FullText[source.position]
	var a = rune(symbol)

	return a, false
}

func GetSourceMock() *SourceMock {
	return &SourceMock{false, "", -1}
}
