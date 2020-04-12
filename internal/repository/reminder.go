package repository

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/sepuka/campaner/internal/domain"
)

const actualBatchLimit = 100

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
		models []domain.Reminder
		err    error
		tx     *pg.Tx
	)

	if tx, err = r.db.Begin(); err != nil {
		return models, nil, err
	}

	err = tx.
		Model(&models).
		Where(`notify_at = ?`, timestamp).
		Limit(actualBatchLimit).
		For(`UPDATE SKIP LOCKED`).
		Select()

	return models, tx, err
}

func (r *ReminderRepository) SetStatus(reminder *domain.Reminder, tx *pg.Tx) (pg.Result, error) {
	return tx.
		Model(reminder).
		Set(`status = ?status`).
		Where(`reminder_id = ?reminder_id`).
		Update()
}

func (r *ReminderRepository) Commit(tx *pg.Tx) error {
	return tx.Commit()
}
