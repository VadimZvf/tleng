package main

import (
	"fmt"
	"os"

	"github.com/VadimZvf/golang/source"
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

type astNode struct {
	nodeType string
	params   []string
	children []astNode
}

func main() {
	filePath := os.Args[1]

	fmt.Println("Reading file: " + filePath)

	file, err := os.Open(filePath)
	check(err)

	var src = source.GetSource(file)
	var buffer = tokenizer_buffer.CreateBuffer(src)
	var tknzr = tokenizer.GetTokenizer(&buffer)

	tokens := tknzr.GetTokens()

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

	// var ast = astNode{token.PROGRAMM.GetCode(), []string{}, []astNode{}}

	// for i := 0; i < len(tokens); i++ {
	// 	t := tokens[i]

	// 	if t.Info.GetCode() == token.VARIABLE_DECLORAION.GetCode() {
	// 		variableNameLocation, isFoundVariableNameLocation := searchFirstKeyWord(tokens, i+1)

	// 		if !isFoundVariableNameLocation {
	// 			panic("variable name decloration INVALID! Cannot find name at positin " + fmt.Sprint(i))
	// 		}

	// 		variableNameToken := tokens[variableNameLocation]

	// 		variableValueLocation, isFoundVariableValue := searchFirstKeyWord(tokens, variableNameLocation+1)

	// 		if !isFoundVariableValue {
	// 			panic("variable value not DEFINED! Cannot find value for valiable " + variableNameToken.Info.GetValue())
	// 		}

	// 		variableValueToken := tokens[variableValueLocation]

	// 		ast.children = append(ast.children, astNode{
	// 			token.VARIABLE_DECLORAION.GetCode(),
	// 			[]string{variableNameToken.Value, variableValueToken.Value},
	// 			[]astNode{},
	// 		})
	// 	}
	// }

	// fmt.Println(ast)
}

// func eatString(tokens []token.Token, startPosition int) (int, bool) {
// 	for i := startPosition; i < len(tokens); i++ {
// 		if tokens[i].Code == token.KEY_WORD.GetCode() {
// 			return i, true
// 		}
// 	}

// 	return 0, false
// }

// func searchFirstKeyWord(tokens []token.Token, startPosition int) (int, bool) {
// 	for i := startPosition; i < len(tokens); i++ {
// 		if tokens[i].Code == token.KEY_WORD.GetCode() {
// 			return i, true
// 		}
// 	}

// 	return 0, false
// }
