package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func NewSpeechIsOverError(pointer int, speech string) error {
	return campanerError{
		errorType:     SpeechIsOverError,
		originalError: errors.New(`pointer is too big`),
		context: map[string]string{
			`pointer`: fmt.Sprint(pointer),
			`speech`:  speech,
		},
	}
}
