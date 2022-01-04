package parser

import (
	"fmt"

	"github.com/VadimZvf/golang/ast"
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error_printer"
	"github.com/VadimZvf/golang/stdout"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_variable_declaration"
	"github.com/VadimZvf/golang/tokenizer"
	"github.com/VadimZvf/golang/tokenizer_buffer"
	"github.com/fatih/color"
)

type iSource interface {
	NextSymbol() (symbol rune, isEnd bool)
}

type Parser struct {
	// Inner props
	source iSource
}

func CreateParser(source iSource) Parser {
	return Parser{
		source: source,
	}
}

func (parser *Parser)Parse(isDebug bool) (*ast_node.ASTNode, error) {
	var buffer = tokenizer_buffer.CreateBuffer(parser.source)
	var tknzr = tokenizer.GetTokenizer(&buffer)

	tokens, parsErr := tknzr.GetTokens()

	if parsErr != nil {
		var std = stdout.CreateStdout()
		parser_error_printer.PrintError(&buffer, &std, parsErr)

		return &ast_node.ASTNode{}, parsErr
	}

	if isDebug {
		for _, v := range tokens {
			color.New(color.FgCyan).Printf(fmt.Sprint(v.StartPosition))
			color.New(color.FgCyan).Printf(" type: ")
			color.New(color.FgYellow).Printf(v.Code)
	
			color.New(color.FgCyan).Printf(" value: \"")
			color.New(color.FgGreen).Printf(v.Value)
	
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

	return ast.CreateAST(tokens)
}