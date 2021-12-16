package ast_factory

import (
	"github.com/VadimZvf/golang/ast_node"
)


type ASTFactory struct {
	root *ast_node.ASTNode
	currentNode *ast_node.ASTNode
}

func CreateASTFactory() ASTFactory {
	var rootNode = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
	}

	return ASTFactory{
		root: &rootNode,
		currentNode: &rootNode,
	}
}

func (factory *ASTFactory) Append(node *ast_node.ASTNode) {
	factory.currentNode.Body = append(factory.currentNode.Body, node)
}

func (factory *ASTFactory) GetCurrent() *ast_node.ASTNode {
	return factory.currentNode
}

func (factory *ASTFactory) GetAST() *ast_node.ASTNode {
	return factory.root
}