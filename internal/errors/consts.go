package errors

const (
	NoType = ErrorType(iota)
	UnConsistentGlossaryError
	TimeIsOverError
	TimeParseError
	NotATimeError
)

type ErrorType uint
