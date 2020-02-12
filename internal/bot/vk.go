package bot

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

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
	var (
		socket   = obj.cfg.Socket
		signals  = make(chan os.Signal, 1)
		stop     = make(chan error, 1)
		listener net.Listener
		err      error
	)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	listener, err = net.Listen(`unix`, socket)
	if err != nil {
		obj.logger.Errorf(`cannot listen to unix socket: %s`, err)
		return err
	}

	defer func() error {
		return os.Remove(socket)
	}()

	go func() {
		<-signals
		if err = listener.Close(); err != nil {
			stop <- errors.Wrap(err, `unable to close HTTP connection`)
		}
	}()

	go obj.server(listener, stop)

	err = <-stop

	return err
}

func (obj *Bot) server(listener net.Listener, c chan<- error) {
	if err := fcgi.Serve(listener, obj); err != nil {
		obj.logger.Errorf(`cannot to serve accept connections: %s`, err)
		c <- err
	}
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