package ast_factory

import (
	"github.com/VadimZvf/golang/ast_node"
)

type ASTFactory struct {
	root      *ast_node.ASTNode
	workStack []*ast_node.ASTNode
}

func CreateASTFactory() ASTFactory {
	var rootNode = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
	}

	return ASTFactory{
		root:      &rootNode,
		workStack: []*ast_node.ASTNode{},
	}
}

func (factory *ASTFactory) Append(node *ast_node.ASTNode) {
	factory.root.Body = append(factory.root.Body, node)
}

func (factory *ASTFactory) PushToWorkStack(node *ast_node.ASTNode) {
	factory.workStack = append(factory.workStack, node)
}

func (factory *ASTFactory) GetLastInWorkStack() *ast_node.ASTNode {
	if len(factory.workStack) == 0 {
		return nil
	}

	return factory.workStack[len(factory.workStack)-1]
}

func (factory *ASTFactory) PopWorkStack() *ast_node.ASTNode {
	var lastItem = factory.GetLastInWorkStack()

	if lastItem == nil {
		return lastItem
	}

	factory.workStack = factory.workStack[:len(factory.workStack)-1]

	return lastItem
}

func (factory *ASTFactory) GetWorkStack() []*ast_node.ASTNode {
	return factory.workStack
}

func (factory *ASTFactory) GetAST() *ast_node.ASTNode {
	return factory.root
}
