package main

import (
	"fmt"
	"os"

	"github.com/VadimZvf/golang/ast_node"
	"github.com/VadimZvf/golang/parser"
	"github.com/VadimZvf/golang/source_file"
)

func main() {
	filePath := os.Args[1]

	var src = source_file.GetSource(filePath)
	var parser = parser.CreateParser(src)

	var astRoot, astError = parser.Parse(true)

	if astError != nil {
		fmt.Println(astError)
	}

	if astError != nil {
		os.Exit(1)
	}

	printASTNode(astRoot, 0, false)
}

func printASTNode(node *ast_node.ASTNode, depth int, isArgument bool) {
	for i := 0; i < depth; i++ {
		if isArgument {
			if i == 0 {
				fmt.Printf("ARG..")
			} else {
				fmt.Printf(".....")
			}
		} else {
			fmt.Printf("|    ")
		}
	}

	fmt.Printf("Code: %s ", node.Code)
	if len(node.Params) > 0 {
		fmt.Printf("Params: ")
		for _, param := range node.Params {
			fmt.Printf("%s=\"%s\" ", param.Name, param.Value)
		}
	}
	fmt.Printf("\n")

	if len(node.Body) > 0 {
		for _, child := range node.Body {
			printASTNode(child, depth+1, false)
		}
	}

	if len(node.Arguments) > 0 {
		for _, child := range node.Arguments {
			printASTNode(child, depth+1, true)
		}
	}
}
