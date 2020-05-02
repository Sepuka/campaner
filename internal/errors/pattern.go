package errors

func NewPatternLengthIncorrect() error {
	return calendarError{
		errorType: PatternLengthIncorrect,
	}
}
