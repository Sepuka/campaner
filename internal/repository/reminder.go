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

func (r *ReminderRepository) Persist(reminder *domain.Reminder) error {
	return r.
		db.
		Insert(reminder)
}

func (r *ReminderRepository) FindActual(timestamp time.Time) ([]domain.Reminder, error) {
	var (
		models []domain.Reminder
		err    error
	)

	err = r.
		db.
		Model(&models).
		Where(`notify_at = ?`, timestamp).
		Limit(actualBatchLimit).
		Select()

	return models, err
}
