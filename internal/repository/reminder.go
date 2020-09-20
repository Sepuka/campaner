package repository

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/sepuka/campaner/internal/domain"
)

const (
	actualBatchLimit = 100
	batchPeriodSec   = 60
)

type ReminderRepository struct {
	db *pg.DB
}

func NewReminderRepository(db *pg.DB) *ReminderRepository {
	return &ReminderRepository{db: db}
}

func (r *ReminderRepository) Add(reminder *domain.Reminder) error {
	return r.
		db.
		Insert(reminder)
}

func (r *ReminderRepository) FindActual(timestamp time.Time) ([]domain.Reminder, *pg.Tx, error) {
	var (
		models      []domain.Reminder
		err         error
		tx          *pg.Tx
		startPeriod = timestamp.Add(-batchPeriodSec * time.Second)
		endPeriod   = timestamp
	)

	if tx, err = r.db.Begin(); err != nil {
		return models, nil, err
	}

	err = tx.
		Model(&models).
		Where(`(notify_at BETWEEN ? AND ?) AND status = ?`, startPeriod, endPeriod, domain.StatusNew).
		Limit(actualBatchLimit).
		For(`UPDATE SKIP LOCKED`).
		Select()

	if models == nil {
		return models, nil, tx.Rollback()
	}

	return models, tx, err
}

func (r *ReminderRepository) SetStatus(reminder *domain.Reminder, tx *pg.Tx) (pg.Result, error) {
	return tx.
		Model(reminder).
		Column(`status`).
		WherePK().
		Update()
}

func (r *ReminderRepository) Commit(tx *pg.Tx) error {
	return tx.Commit()
}

func (r *ReminderRepository) Scheduled(userId int, limit uint32) ([]domain.Reminder, error) {
	var (
		models []domain.Reminder
		err    error
	)

	err = r.
		db.
		Model(&models).
		Where(`notify_at > ? AND user_id = ?`, time.Now(), userId).
		Limit(int(limit)).
		Select()

	return models, err
}

func (r *ReminderRepository) Cancel(taskId int64, userId int) error {
	var (
		model = &domain.Reminder{
			ReminderId: int(taskId),
		}
		err error
	)

	err = r.
		db.
		Model(model).
		Where(`reminder_id = ? AND user_id = ? AND status = ?`, taskId, userId, domain.StatusNew).
		Select()

	if err != nil {
		return err
	}

	model.Status = domain.StatusCanceled

	return r.
		db.
		Update(model)
}
