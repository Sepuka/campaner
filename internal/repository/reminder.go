package repository

import (
	"github.com/go-pg/pg"
	"github.com/sepuka/campaner/internal/domain"
)

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
