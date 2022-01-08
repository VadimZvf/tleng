package parser_error_printer

import (
	"github.com/VadimZvf/golang/parser_error"
)

type iStdout interface {
	Print(line string)
	PrintError(line string)
}

func PrintError(code string, std iStdout, err error) {
	re, ok := err.(parser_error.ParserError)

	if !ok {
		panic(err)
	}

	var i = 0

	for ; i < re.StartPosition; i++ {
		std.Print(string(code[i]))
	}

	// Print error parth
	for ; i <= re.EndPosition; i++ {
		std.PrintError(string(code[i]))
	}

	// Print valid part, at the same line
	for ; i < len(code) && string(code[i]) != "\n"; i++ {
		std.Print(string(code[i]))
	}

	std.Print("\n")
	std.PrintError(re.Message + "\n")

	for ; i < len(code); i++ {
		std.Print(string(code[i]))
	}

	std.Print("\n")
}
