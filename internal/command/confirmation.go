package command

import (
	"github.com/sepuka/campaner/internal/config"
	"github.com/sepuka/campaner/internal/context"
)

type Confirmation struct {
	cfg config.Server
}

func NewConfirmation(cfg config.Server) *Confirmation {
	return &Confirmation{
		cfg: cfg,
	}
}

func (o *Confirmation) Exec(req *context.Request, resp *Result) error {
	resp.SetResponse(o.cfg.Confirmation)

	return nil
}

func (o *Confirmation) Precept() []string {
	return []string{
		`remind`,
		`напомни`,
	}
}
