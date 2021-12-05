package parser_error_printer

import (
	"github.com/VadimZvf/golang/parser_error"
)

type iBuffer interface {
	GetReadedCode() (value string)
}

type iStdout interface {
	PrintLine(line string)
	PrintSymbol(symbol string)
	SetErrorColor()
	SetDefaultColor()
}

func PrintError(buffer iBuffer, std iStdout, err error) {
	re, ok := err.(parser_error.ParserError)

	if !ok {
		panic(err)
	}

	var i = 0
	var code = buffer.GetReadedCode()

	for ; i < re.StartPosition; i++ {
		std.PrintSymbol(string(code[i]))
	}

	// Print error parth
	std.SetErrorColor()
	for ; i <= re.EndPosition; i++ {
		std.PrintSymbol(string(code[i]))
	}
	std.SetDefaultColor()

	// Print valid part, at the same line
	for ; i < len(code) && string(code[i]) != "\n"; i++ {
		std.PrintSymbol(string(code[i]))
	}

	std.SetErrorColor()
	std.PrintLine(re.Message)
	std.SetDefaultColor()

	for ; i < len(code); i++ {
		std.PrintSymbol(string(code[i]))
	}

	std.PrintLine("")
}
