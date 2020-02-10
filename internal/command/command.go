package command

import (
	"net/http"

	"github.com/sepuka/campaner/internal/context"
)

type HandlerMap map[string]Executor

type Executor interface {
	Exec(*context.Request, http.ResponseWriter) error
}

type Preceptable interface {
	Precept() []string
}
