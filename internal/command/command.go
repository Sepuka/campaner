package command

import "github.com/sepuka/campaner/internal/context"

type HandlerMap map[string]Executor

type Result struct {
	response string
}

func (o *Result) SetResponse(response string) {
	o.response = response
}

type Executor interface {
	Exec(*context.Request, *Result) error
}

type Preceptable interface {
	Precept() []string
}
