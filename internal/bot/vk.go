package bot

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	"github.com/sepuka/campaner/internal/config"

	"github.com/sepuka/campaner/internal/middleware"

	"github.com/sepuka/campaner/internal/command"
	"go.uber.org/zap"
)

type Bot struct {
	commands command.HandlerMap
	logger   *zap.SugaredLogger
	handler  middleware.HandlerFunc
	cfg      config.Server
}

func NewBot(
	logger *zap.SugaredLogger,
	commandsMap command.HandlerMap,
	handler middleware.HandlerFunc,
	cfg config.Server,
) *Bot {
	return &Bot{
		commands: commandsMap,
		logger:   logger,
		handler:  handler,
		cfg:      cfg,
	}
}

func (o *Bot) Listen() error {
	socket := o.cfg.Socket
	defer os.Remove(socket)

	l, err := net.Listen(`unix`, socket)
	if err != nil {
		return err
	}
	fs := new(FCgiServer)

	return fcgi.Serve(l, fs)
}

type (
	FCgiServer struct{}
)

func (this *FCgiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`500 Server error`))
		}
	}()
	w.Write([]byte(`OK Server<br>text`))
}
