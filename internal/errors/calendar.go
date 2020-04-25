package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type calendarError struct {
	errorType     ErrorType
	originalError error
	context       map[string]string
}

func (e calendarError) Error() string {
	return e.originalError.Error()
}

func (e calendarError) New(msg string) error {
	return calendarError{
		errorType:     e.errorType,
		originalError: errors.New(msg),
	}
}

func (e calendarError) Wrap(err error, msg string) error {
	return e.Wrapf(err, msg)
}

func (e calendarError) Wrapf(err error, msg string, args ...interface{}) error {
	wrappedErr := errors.Wrapf(err, msg, args...)

	return calendarError{
		errorType:     e.errorType,
		originalError: wrappedErr,
	}
}

func NewNotATimeError() error {
	return calendarError{
		errorType:     NotATimeError,
		originalError: errors.New(`there is not any info about time`),
	}
}

func NewUnConsistentGlossaryError(word string, glossary []string) error {
	return calendarError{
		errorType:     UnConsistentGlossaryError,
		originalError: errors.New(`got unknown keyword`),
		context: map[string]string{
			`got`: word,
			`can`: strings.Join(glossary, `,`),
		},
	}
}

func NewUnknownDimensionError(value float64, dimension string) error {
	return calendarError{
		errorType:     UnknownDimensionError,
		originalError: errors.New(`got unknown time dimension`),
		context: map[string]string{
			`value`:     fmt.Sprint(value),
			`dimension`: dimension,
		},
	}
}

func NewInvalidTimeValueError(value int) error {
	return calendarError{
		errorType:     InvalidTimeError,
		originalError: errors.New(`got invalid time`),
		context: map[string]string{
			`value`: fmt.Sprint(value),
		},
	}
}
