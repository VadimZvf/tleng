package source_file

import (
	"fmt"
	"os"
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

func (source Source) NextSymbol() (symbol string, isEnd bool) {
	n1, err := source.file.Read(source.buffer)

	return string(source.buffer[:n1]), err != nil
}
