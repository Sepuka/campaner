package errors

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func NewUnknownBrokerError(status int, err string) error {
	return campanerError{
		errorType:     Tasks,
		originalError: errors.New(err),
		context: map[string]string{
			`status`: strconv.Itoa(status),
		},
	}
}

func NewWrongUserError(owner, executor int) error {
	return campanerError{
		errorType:     Tasks,
		originalError: NewSecurityError(),
		context: map[string]string{
			`owner`:    strconv.Itoa(owner),
			`executor`: strconv.Itoa(executor),
		},
	}
}

func NewShiftError(duration time.Duration) error {
	return campanerError{
		errorType: Tasks,
		context: map[string]string{
			`duration`: duration.String(),
		},
	}
}

func NewWrongStatusError(actualStatus, wantStatus int) error {
	return campanerError{
		errorType: Tasks,
		context: map[string]string{
			`actual`: strconv.Itoa(actualStatus),
			`want`:   strconv.Itoa(wantStatus),
		},
	}
}
