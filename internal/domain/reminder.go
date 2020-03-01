package domain

import (
	"context"
	"errors"
	"time"

	"github.com/go-pg/pg/orm"
)

const toShortTime = 5 * time.Second

type (
	ReminderRepository interface {
		Persist(reminder *Reminder) error
	}

	Reminder struct {
		ReminderId int       `pg:",pk"`
		CreatedAt  time.Time `pg:"default:now()"`
		Whom       int       `sql:"user_id"`
		What       string    `sql:"content",pg:"notnull"`
		NotifyAt   time.Time
		When       time.Duration `sql:"-"`
	}
)

func NewReminder(whom int, what string, when time.Duration) *Reminder {
	return &Reminder{
		Whom: whom,
		What: what,
		When: when,
	}
}

func NewImmediateReminder(
	whom int,
	what string,
) *Reminder {
	return &Reminder{
		When: time.Nanosecond,
		Whom: whom,
		What: what,
	}
}

func (r *Reminder) BeforeInsert(ctx context.Context, db orm.DB) error {
	if r.When == 0 {
		return errors.New(`"when" column is not initialized`)
	}

	r.NotifyAt = time.Now().Add(r.When)
	r.CreatedAt = time.Now()

	return nil
}

func (r *Reminder) IsImmediate() bool {
	return r.When < toShortTime
}

func (r *Reminder) IsValid() bool {
	return r.When > 0
}
