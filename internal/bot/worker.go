package bot

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sepuka/campaner/internal/api/method"

	"github.com/go-pg/pg"

	"go.uber.org/zap"

	"github.com/sepuka/campaner/internal/domain"
)

type Worker struct {
	repo   domain.ReminderRepository
	api    *method.MessagesSend
	logger *zap.SugaredLogger
}

func NewWorker(
	repo domain.ReminderRepository,
	logger *zap.SugaredLogger,
	api *method.MessagesSend,
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
		tx   *pg.Tx
	)

	if rows, tx, err = w.repo.FindActual(moment); err != nil {
		w.
			logger.
			With(
				zap.Error(err),
			).
			Error(`error finding actual reminders`)
		return
	}

	if rows == nil {
		return
	}

	for _, row = range rows {
		w.remind(&row, tx)
	}

	if err = w.repo.Commit(tx); err != nil {
		w.
			logger.
			With(
				zap.Error(err),
			).
			Errorf(`error while commit reminder statuses`)
	}
}

func (w *Worker) remind(reminder *domain.Reminder, tx *pg.Tx) {
	var (
		err    error
		status = domain.StatusSuccess
	)

	if err = w.api.Send(reminder.Whom, reminder.What); err != nil {
		status = domain.StatusFailed
		w.
			logger.
			With(
				zap.String(`text`, reminder.What),
				zap.Error(err),
			).
			Error(`send api message error`)
	}

	reminder.Status = status
	if _, err = w.repo.SetStatus(reminder, tx); err != nil {
		w.
			logger.
			With(
				zap.Any(`reminder`, reminder),
				zap.Error(err),
			).
			Error(`error updating reminder status`)
	}
}
