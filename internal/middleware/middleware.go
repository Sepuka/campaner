package middleware

import (
	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/context"
)

type HandlerFunc func(command.Executor, *context.Request, *command.Result) error

func final(handler command.Executor, req *context.Request, resp *command.Result) error {
	return handler.Exec(req, resp)
}

func BuildHandlerChain(handlers []func(HandlerFunc) HandlerFunc) HandlerFunc {
	if len(handlers) == 0 {
		return final
	}

	return handlers[0](BuildHandlerChain(handlers[1:]))
}
