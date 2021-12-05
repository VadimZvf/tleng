package main

import (
	"fmt"
	"os"

	"github.com/VadimZvf/golang/parser_error_printer"
	"github.com/VadimZvf/golang/source_file"
	"github.com/VadimZvf/golang/stdout"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_variable_declaration"
	"github.com/VadimZvf/golang/tokenizer"
	"github.com/VadimZvf/golang/tokenizer_buffer"

	"github.com/fatih/color"
)

func main() {
	filePath := os.Args[1]

	var src = source_file.GetSource(filePath)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tknzr = tokenizer.GetTokenizer(&buffer)

	tokens, parsErr := tknzr.GetTokens()

	if parsErr != nil {
		var std = stdout.CreateStdout()
		parser_error_printer.PrintError(&buffer, &std, parsErr)

		os.Exit(1)
	}

	for _, v := range tokens {
		color.New(color.FgCyan).Printf(fmt.Sprint(v.StartPosition))
		color.New(color.FgCyan).Printf(" type: ")
		color.New(color.FgYellow).Printf(v.Code)

		color.New(color.FgCyan).Printf(" value: \"")
		color.New(color.FgGreen).Printf(v.DebugValue)

		if v.Code == token_function_declaration.FUNCTION_DECLARATION || v.Code == token_variable_declaration.VARIABLE_DECLARAION {
			for _, param := range v.Params {
				color.New(color.FgGreen).Printf(param.Name)
				color.New(color.FgGreen).Printf("=")
				color.New(color.FgGreen).Printf(param.Value)
			}
		}

		color.New(color.FgCyan).Printf("\"\n")
	}
}
