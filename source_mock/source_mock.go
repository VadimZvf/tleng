package source_mock

type SimpleSourceMock struct {
	IsEnd           bool
	NextSymbolValue string
}

func (source *SimpleSourceMock) NextSymbol() (symbol string, isEnd bool) {
	return source.NextSymbolValue, source.IsEnd
}

func GetSimpleSource() *SimpleSourceMock {
	return &SimpleSourceMock{false, "a"}
}

type SourceMock struct {
	IsEnd    bool
	FullText string
	position int
}

func (source *SourceMock) NextSymbol() (symbol string, isEnd bool) {
	source.position += 1

	if source.position >= len(source.FullText) {
		return "", true
	}

	return string(source.FullText[source.position]), false
}

func GetSourceMock() *SourceMock {
	return &SourceMock{false, "", -1}
}
