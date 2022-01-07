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

func MergeRuntimeErrors(first error, second error) error {
	firstError, firstCastOk := first.(RuntimeError)

	if !firstCastOk {
		return firstError
	}

	secondError, secondCastOk := second.(RuntimeError)

	if !secondCastOk {
		return firstError
	}

	firstError.Message = secondError.Message + "\n  " + firstError.Message
	firstError.StartPosition = secondError.StartPosition
	firstError.EndPosition = secondError.EndPosition

	return firstError
}
