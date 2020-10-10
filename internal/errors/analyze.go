package errors

import (
	"github.com/sepuka/campaner/internal/context"
)

func NewInvalidSpeechPayloadFormatError(msg context.Message, err error) error {
	return campanerError{
		errorType:     InvalidSpeechPayload,
		originalError: err,
		context: map[string]string{
			`text`:    msg.Text,
			`payload`: msg.Payload,
		},
	}
}

func NewInvalidSpeechPayloadButtonError(msg context.Message, err error) error {
	return campanerError{
		errorType:     InvalidSpeechPayload,
		originalError: err,
		context: map[string]string{
			`text`:    msg.Text,
			`payload`: msg.Payload,
		},
	}
}

func NewTaskManagerError(msg string) error {
	return campanerError{
		errorType: Storage,
		context: map[string]string{
			`text`: msg,
		},
	}
}
