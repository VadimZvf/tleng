package parser

import (
	"fmt"

	"github.com/VadimZvf/golang/ast"
	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/token_function_declaration"
	"github.com/VadimZvf/golang/token_variable_declaration"
	"github.com/VadimZvf/golang/tokenizer"
	"github.com/VadimZvf/golang/tokenizer_buffer"
)

type ISource interface {
	NextSymbol() (symbol rune, isEnd bool)
}

type iStdout interface {
	Print(line string)
	PrintError(line string)
}

type Parser struct {
	// Inner props
	source ISource
	stdout iStdout
	buffer tokenizer_buffer.Buffer
}

func CreateParser(source ISource, stdout iStdout) Parser {
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
		return &ast_node.ASTNode{}, parsErr
	}

	if isDebug {
		parser.stdout.Print("__________________TOKENS______________________\n")
		for _, v := range tokens {
			parser.stdout.Print(fmt.Sprint(v.StartPosition))
			parser.stdout.Print(" type: ")
			parser.stdout.PrintError(v.Code)

			parser.stdout.Print(" value: \"")
			parser.stdout.Print(v.Value)

			if v.Code == token_function_declaration.FUNCTION_DECLARATION || v.Code == token_variable_declaration.VARIABLE_DECLARAION {
				for _, param := range v.Params {
					parser.stdout.PrintError(param.Name)
					parser.stdout.PrintError("=")
					parser.stdout.PrintError(param.Value)
					parser.stdout.PrintError(" ")
				}
			}

			parser.stdout.Print("\"\n")
		}
		parser.stdout.Print("_____________________________________________\n")
	}

	var ast, astError = ast.CreateAST(tokens)

	if isDebug && ast != nil {
		parser.stdout.Print("___________________AST_______________________\n")
		printASTNode(parser.stdout, ast, 0, false)
		parser.stdout.Print("_____________________________________________\n")
	}

	return ast, astError
}

func (parser *Parser) GetSourceCode() string {
	return parser.buffer.GetReadedCode()
}

func printASTNode(stdout iStdout, node *ast_node.ASTNode, depth int, isArgument bool) {
	for i := 0; i < depth; i++ {
		if isArgument {
			if i == 0 {
				stdout.Print("ARG..")
			} else {
				stdout.Print(".....")
			}
		} else {
			stdout.Print("|    ")
		}
	}

	stdout.Print(fmt.Sprintf("Code: %s ", node.Code))
	if len(node.Params) > 0 {
		stdout.Print("Params: ")
		for _, param := range node.Params {
			stdout.Print(fmt.Sprintf("%s=\"%s\" ", param.Name, param.Value))
		}
	}
	stdout.Print("\n")

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
