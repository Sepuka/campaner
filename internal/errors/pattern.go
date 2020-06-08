package errors

func NewPatternLengthIncorrect() error {
	return campanerError{
		errorType: PatternLengthIncorrect,
	}
}
