package errors

const (
	NoType = ErrorType(iota)
	SecurityError
	UnConsistentGlossaryError
	ItIsPastTimeError
	UnknownDimensionError
	NotATimeError
	InvalidTimeError

	SpeechIsOverError

	PatternLengthIncorrect
	InvalidSpeechPayload
	Storage
	Tasks
)

type ErrorType uint

func GetType(err error) ErrorType {
	if errType, ok := err.(campanerError); ok {
		return errType.errorType
	}

	return NoType
}

func IsNotATimeError(err error) bool {
	return GetType(err) == NotATimeError
}

func IsSpeechIsOverError(err error) bool {
	return GetType(err) == SpeechIsOverError
}
