package parser

import (
	"fmt"

	"github.com/VadimZvf/golang/ast"
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser_error_printer"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_variable_declaration"
	"github.com/VadimZvf/golang/tokenizer"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

type iSource interface {
	NextSymbol() (symbol rune, isEnd bool)
}

type iStdout interface {
	PrintLine(line string)
	PrintSymbol(symbol string)
	SetErrorColor()
	SetDefaultColor()
}

type Parser struct {
	// Inner props
	source iSource
	stdout iStdout
	buffer tokenizer_buffer.Buffer
}

func CreateParser(source iSource, stdout iStdout) Parser {
	var buffer = tokenizer_buffer.CreateBuffer(source)

	return Parser{
		source: source,
		stdout: stdout,
		buffer: buffer,
	}
}

func (parser *Parser) Parse(isDebug bool) (*ast_node.ASTNode, error) {
	var tknzr = tokenizer.GetTokenizer(&parser.buffer)

	tokens, parsErr := tknzr.GetTokens()

	if parsErr != nil {
		parser_error_printer.PrintError(&parser.buffer, parser.stdout, parsErr)

		return &ast_node.ASTNode{}, parsErr
	}

	if isDebug {
		for _, v := range tokens {
			parser.stdout.SetDefaultColor()
			parser.stdout.PrintSymbol(fmt.Sprint(v.StartPosition))
			parser.stdout.PrintSymbol(" type: ")
			parser.stdout.SetErrorColor()
			parser.stdout.PrintSymbol(v.Code)
			parser.stdout.SetDefaultColor()

			parser.stdout.PrintSymbol(" value: \"")
			parser.stdout.PrintSymbol(v.Value)

			parser.stdout.SetErrorColor()
			if v.Code == token_function_declaration.FUNCTION_DECLARATION || v.Code == token_variable_declaration.VARIABLE_DECLARAION {
				for _, param := range v.Params {
					parser.stdout.PrintSymbol(param.Name)
					parser.stdout.PrintSymbol("=")
					parser.stdout.PrintSymbol(param.Value)
					parser.stdout.PrintSymbol(" ")
				}
			}
			parser.stdout.SetDefaultColor()

			parser.stdout.PrintSymbol("\"\n")
		}
	}

	var ast, astError = ast.CreateAST(tokens)

	if astError != nil && isDebug {
		parser.stdout.PrintLine("_____________________________________________")
		parser_error_printer.PrintError(&parser.buffer, parser.stdout, astError)
		parser.stdout.PrintLine("_____________________________________________")
	}

	if isDebug && ast != nil {
		printASTNode(parser.stdout, ast, 0, false)
	}

	return ast, astError
}

func (parser *Parser) GetSourceCode() string {
	return parser.buffer.GetReadedCode()
}

func printASTNode(stdout iStdout, node *ast_node.ASTNode, depth int, isArgument bool) {
	if isArgument {
		stdout.SetErrorColor()
	} else {
		stdout.SetDefaultColor()
	}

	for i := 0; i < depth; i++ {
		if isArgument {
			if i == 0 {
				stdout.PrintSymbol("ARG..")
			} else {
				stdout.PrintSymbol(".....")
			}
		} else {
			stdout.PrintSymbol("|    ")
		}
	}

	stdout.PrintSymbol(fmt.Sprintf("Code: %s ", node.Code))
	if len(node.Params) > 0 {
		stdout.PrintSymbol("Params: ")
		for _, param := range node.Params {
			stdout.PrintSymbol(fmt.Sprintf("%s=\"%s\" ", param.Name, param.Value))
		}
	}
	stdout.PrintSymbol("\n")

	if len(node.Body) > 0 {
		for _, child := range node.Body {
			printASTNode(stdout, child, depth+1, false)
		}
	}

	if len(node.Arguments) > 0 {
		for _, child := range node.Arguments {
			printASTNode(stdout, child, depth+1, true)
		}
	}
}
