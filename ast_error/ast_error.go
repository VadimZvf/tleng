package ast_error

type AstError struct {
	Message       string
	StartPosition int
	EndPosition   int
}

func (err AstError) Error() string {
	return err.Message
}

func CreateError(message string, start int, end int) AstError {
	return AstError{
		Message:       message,
		StartPosition: start,
		EndPosition:   end,
	}
}
