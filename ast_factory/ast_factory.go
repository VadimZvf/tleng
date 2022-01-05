package ast_factory

import (
	"github.com/VadimZvf/golang/ast_node"
)

type ASTFactory struct {
	root *ast_node.ASTNode
}

func CreateASTFactory() ASTFactory {
	var rootNode = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
	}

	return ASTFactory{
		root: &rootNode,
	}
}

func (factory *ASTFactory) Append(node *ast_node.ASTNode) {
	factory.root.Body = append(factory.root.Body, node)
}

func (factory *ASTFactory) GetAST() *ast_node.ASTNode {
	return factory.root
}
