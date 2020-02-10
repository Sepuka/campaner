package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sepuka/campaner/internal/command"
	"github.com/sepuka/campaner/internal/context"
)

func Panic(next HandlerFunc) HandlerFunc {
	return func(exec command.Executor, req *context.Request, writer http.ResponseWriter) error {
		var err error

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %s\n"+
					"command `%s`\n"+
					"stacktrace from panic: %s\n",
					r, req.Type, string(debug.Stack()))
				err = errors.New(`internal error`)
			}
		}()

		err = next(exec, req, writer)

		return err
	}
}
