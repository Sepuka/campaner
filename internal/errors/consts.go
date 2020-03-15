package errors

const (
	NoType = ErrorType(iota)
	UnConsistentGlossaryError
	ItIsPastTimeError
	UnknownDimensionError
	NotATimeError
	InvalidTimeError
)

type ErrorType uint
