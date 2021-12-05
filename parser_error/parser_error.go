package parser_error

type ParserError struct {
	Message       string
	StartPosition int
	EndPosition   int
}

func (err ParserError) Error() string {
	return err.Message
}

func CreateError(message string, start int, end int) ParserError {
	return ParserError{
		Message:       message,
		StartPosition: start,
		EndPosition:   end,
	}
}
