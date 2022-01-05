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

	printASTNode(astRoot, 0)
}

func printASTNode(node *ast_node.ASTNode, depth int) {
	for i := 0; i < depth; i++ {
		fmt.Printf("    ")
	}
	fmt.Printf("Code: %s ", node.Code)
	fmt.Println("Params:", node.Params)

	if len(node.Body) > 0 {
		for _, child := range node.Body {
			printASTNode(child, depth+1)
		}
	}

	if len(node.Arguments) > 0 {
		for _, child := range node.Body {
			printASTNode(child, depth+1)
		}
	}
}
