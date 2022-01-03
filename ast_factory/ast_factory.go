package ast_factory

import (
	"github.com/VadimZvf/golang/ast_node"
)

type ASTFactory struct {
	root        *ast_node.ASTNode
	currentNode *ast_node.ASTNode
	path        []*ast_node.ASTNode

	workInProgress *ast_node.ASTNode
}

func CreateASTFactory() ASTFactory {
	var rootNode = ast_node.ASTNode{
		Code: ast_node.AST_NODE_CODE_ROOT,
	}

	return ASTFactory{
		root:        &rootNode,
		currentNode: &rootNode,
		path:        []*ast_node.ASTNode{},
	}
}

func (factory *ASTFactory) Append(node *ast_node.ASTNode) {
	factory.currentNode.Body = append(factory.currentNode.Body, node)
}

func (factory *ASTFactory) GetCurrent() *ast_node.ASTNode {
	return factory.currentNode
}

func (factory *ASTFactory) MovePointerToParent() {
	if len(factory.path) > 0 {
		var parentNode = factory.path[len(factory.path)-1]
		factory.currentNode = parentNode
		factory.path = factory.path[:len(factory.path)-1]
	}
}

func (factory *ASTFactory) MovePointerLastNodeBody() {
	var bodyItemsCount = len(factory.currentNode.Body)
	var lastItemInBody = factory.currentNode.Body[bodyItemsCount-1]
	factory.path = append(factory.path, factory.currentNode)
	factory.currentNode = lastItemInBody
}

func (factory *ASTFactory) GetAST() *ast_node.ASTNode {
	return factory.root
}
