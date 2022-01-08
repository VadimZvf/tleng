package source_string

type SourceString struct {
	IsEnd    bool
	FullText string
	position int
}

func (source *SourceString) NextSymbol() (rune, bool) {
	source.position += 1

	if source.position >= len(source.FullText) {
		return rune(0), true
	}

	var symbol = source.FullText[source.position]
	var a = rune(symbol)

	return a, false
}

func GetSource(codeText string) *SourceString {
	return &SourceString{false, codeText, -1}
}
