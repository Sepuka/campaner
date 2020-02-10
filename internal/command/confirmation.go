package command

import (
	"net/http"

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

func (o *Confirmation) Exec(req *context.Request, resp http.ResponseWriter) error {
	_, err := resp.Write([]byte(o.cfg.Confirmation))

	return err
}

func (o *Confirmation) Precept() []string {
	return []string{
		`confirmation`,
	}
}
