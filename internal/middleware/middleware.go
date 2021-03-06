package middleware

import (
	"net/http"

	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/context"
)

type HandlerFunc func(command.Executor, *context.Request, http.ResponseWriter) error

func final(handler command.Executor, req *context.Request, resp http.ResponseWriter) error {
	return handler.Exec(req, resp)
}

func BuildHandlerChain(handlers []func(HandlerFunc) HandlerFunc) HandlerFunc {
	if len(handlers) == 0 {
		return final
	}

	return handlers[0](BuildHandlerChain(handlers[1:]))
}
