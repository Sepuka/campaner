package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type campanerError struct {
	errorType     ErrorType
	originalError error
	context       map[string]string
}

func (e campanerError) Error() string {
	if e.originalError != nil {
		return e.originalError.Error()
	}

	return fmt.Sprintf(`error code %d`, e.errorType)
}

func (e campanerError) New(msg string) error {
	return campanerError{
		errorType:     e.errorType,
		originalError: errors.New(msg),
	}
}

func (e campanerError) Wrap(err error, msg string) error {
	return e.Wrapf(err, msg)
}

func (e campanerError) Wrapf(err error, msg string, args ...interface{}) error {
	wrappedErr := errors.Wrapf(err, msg, args...)

	return campanerError{
		errorType:     e.errorType,
		originalError: wrappedErr,
	}
}

func NewNotATimeError() error {
	return campanerError{
		errorType:     NotATimeError,
		originalError: errors.New(`there is not any info about time`),
	}
}

func NewUnConsistentGlossaryError(word string, glossary []string) error {
	return campanerError{
		errorType:     UnConsistentGlossaryError,
		originalError: errors.New(`got unknown keyword`),
		context: map[string]string{
			`got`: word,
			`can`: strings.Join(glossary, `,`),
		},
	}
}

func NewUnknownDimensionError(value float64, dimension string) error {
	return campanerError{
		errorType:     UnknownDimensionError,
		originalError: errors.New(`got unknown time dimension`),
		context: map[string]string{
			`value`:     fmt.Sprint(value),
			`dimension`: dimension,
		},
	}
}

func NewInvalidTimeValueError(value int) error {
	return campanerError{
		errorType:     InvalidTimeError,
		originalError: errors.New(`got invalid time`),
		context: map[string]string{
			`value`: fmt.Sprint(value),
		},
	}
}
