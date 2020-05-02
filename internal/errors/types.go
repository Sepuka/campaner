package errors

const (
	NoType = ErrorType(iota)
	UnConsistentGlossaryError
	ItIsPastTimeError
	UnknownDimensionError
	NotATimeError
	InvalidTimeError

	SpeechIsEmptyError
	SpeechIsOverError

	PatternLengthIncorrect
)

type ErrorType uint

func GetType(err error) ErrorType {
	if errType, ok := err.(calendarError); ok {
		return errType.errorType
	}

	return NoType
}

func IsNotATimeError(err error) bool {
	return GetType(err) == NotATimeError
}
