package bot

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sepuka/campaner/internal/api"
	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/domain"
)

type Worker struct {
	repo   domain.ReminderRepository
	api    *api.SendMessage
	logger *zap.SugaredLogger
}

func NewWorker(
	repo domain.ReminderRepository,
	logger *zap.SugaredLogger,
	api *api.SendMessage,
) *Worker {
	return &Worker{
		repo:   repo,
		logger: logger,
		api:    api,
	}
}

func (w *Worker) Work() error {
	var (
		stop    bool
		signals = make(chan os.Signal, 1)
		moment  time.Time
	)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		stop = true
	}()

	for !stop {
		moment = time.Now().Round(time.Second).UTC()
		go w.notify(moment)
		time.Sleep(time.Second)
	}

	return nil
}

func (w *Worker) notify(moment time.Time) {
	var (
		rows []domain.Reminder
		row  domain.Reminder
		err  error
	)

	if rows, err = w.repo.FindActual(moment); err != nil {
		w.
			logger.
			With(
				zap.Error(err),
			).
			Error(`error finding actual reminders`)
		return
	}

	for _, row = range rows {
		go w.remind(row)
	}
}

func (w *Worker) remind(reminder domain.Reminder) {
	var (
		err error
	)

	if err = w.api.Send(reminder.Whom, reminder.What); err != nil {
		w.
			logger.
			With(
				zap.String(`text`, reminder.What),
				zap.Error(err),
			).
			Error(`send api message error`)
	}
}
