package runtime_error

import "github.com/VadimZvf/golang/ast_node"

type RuntimeError struct {
	Message       string
	StartPosition int
	EndPosition   int
}

func (err RuntimeError) Error() string {
	return err.Message
}

func CreateError(message string, node *ast_node.ASTNode) RuntimeError {
	return RuntimeError{
		Message:       message,
		StartPosition: node.StartPosition,
		EndPosition:   node.EndPosition,
	}
}
