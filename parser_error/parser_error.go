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

func MergeParserErrors(first error, second error) error {
	firstParserError, firstCastOk := first.(ParserError)

	if !firstCastOk {
		return firstParserError
	}

	secondParserError, secondCastOk := second.(ParserError)

	if !secondCastOk {
		return firstParserError
	}

	firstParserError.Message = secondParserError.Message + "\n  " + firstParserError.Message
	firstParserError.StartPosition = secondParserError.StartPosition
	firstParserError.EndPosition = secondParserError.EndPosition

	return firstParserError
}
