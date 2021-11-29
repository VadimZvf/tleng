package main

import (
	"fmt"
	"os"

	"github.com/VadimZvf/golang/parser_error_printer"
	"github.com/VadimZvf/golang/source_file"
	"github.com/VadimZvf/golang/stdout"
	"github.com/VadimZvf/golang/token"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_variable_decloration"
	"github.com/VadimZvf/golang/tokenizer"
	"github.com/VadimZvf/golang/tokenizer_buffer"

	"github.com/fatih/color"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	filePath := os.Args[1]

	fmt.Println("Reading file: " + filePath)

	file, err := os.Open(filePath)
	check(err)

	var src = source_file.GetSource(file)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tknzr = tokenizer.GetTokenizer(&buffer)

	tokens, parsErr := tknzr.GetTokens()

	if parsErr != nil {
		var std = stdout.CreateStdout()
		parser_error_printer.PrintError(&buffer, &std, parsErr)

		os.Exit(1)
	}

	file.Close()

	for _, v := range tokens {
		color.New(color.FgCyan).Printf(fmt.Sprint(v.Position))
		color.New(color.FgCyan).Printf(" type: ")
		color.New(color.FgYellow).Printf(v.Code)

		color.New(color.FgCyan).Printf(" value: \"")

		if v.Code != token.NEW_LINE {
			color.New(color.FgGreen).Printf(v.DebugValue)
		}

		if v.Code == token_function_declaration.FUNCTION_DECLARATION || v.Code == token_variable_decloration.VARIABLE_DECLARAION {
			for _, param := range v.Params {
				color.New(color.FgGreen).Printf(param.Name)
				color.New(color.FgGreen).Printf("=")
				color.New(color.FgGreen).Printf(param.Value)
			}
		}

		color.New(color.FgCyan).Printf("\"\n")
	}
}
