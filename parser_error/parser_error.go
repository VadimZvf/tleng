package parser_error

type ParserError struct {
	Message  string
	Position int
	Length   int
}

func (err ParserError) Error() string {
	return err.Message
}

func CreateError(message string, position int, length int) ParserError {
	return ParserError{
		Message:  message,
		Position: position,
		Length:   length,
	}
}
