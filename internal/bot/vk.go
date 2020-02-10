package bot

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	"github.com/sepuka/campaner/internal/context"

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

func (obj *Bot) Listen() error {
	socket := obj.cfg.Socket
	defer os.Remove(socket)

	l, err := net.Listen(`unix`, socket)
	if err != nil {
		return err
	}

	return fcgi.Serve(l, obj)
}

func (obj *Bot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		callback = &context.Request{}
		decoder  = json.NewDecoder(r.Body)
		err      error
	)

	defer r.Body.Close()
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`500 Server error`))
		}
	}()

	obj.logger.Info(`incoming request`)

	if err = decoder.Decode(callback); err != nil {
		if _, err = w.Write([]byte(`invalid json`)); err != nil {
			obj.logger.Errorf(`cannot write error message about invalid incoming json %s`, err)
		}
		w.WriteHeader(400)

		return
	}

	if finalHandler, ok := obj.commands[callback.Type]; ok {
		obj.handler(finalHandler, callback, w)
	} else {
		if _, err = w.Write([]byte(`unknown type field`)); err != nil {
			obj.logger.Errorf(`cannot write error message about unknown type field %s`, err)
		}
		w.WriteHeader(400)
	}
}
