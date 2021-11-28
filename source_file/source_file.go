package source_file

import "os"

type Source struct {
	file   *os.File
	buffer []byte
}

func GetSource(file *os.File) Source {
	return Source{
		file,
		make([]byte, 1),
	}
}

func (source Source) NextSymbol() (symbol string, isEnd bool) {
	n1, err := source.file.Read(source.buffer)

	return string(source.buffer[:n1]), err != nil
}
