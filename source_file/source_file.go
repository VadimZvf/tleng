package source_file

import (
	"fmt"
	"os"
	"unicode/utf8"
)

type Source struct {
	file   *os.File
	buffer []byte
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetSource(filePath string) Source {
	fmt.Println("Reading file: " + filePath)

	file, err := os.Open(filePath)
	check(err)

	return Source{
		file,
		make([]byte, 1),
	}
}

func (source *Source) Close() {
	source.file.Close()
}

func (source Source) NextSymbol() (rune, bool) {
	n1, err := source.file.Read(source.buffer)

	symbol, _ := utf8.DecodeRune(source.buffer[:n1])

	return symbol, err != nil
}
