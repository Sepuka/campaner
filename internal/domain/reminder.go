package domain

import (
	"context"
	"errors"
	"time"

	"github.com/go-pg/pg"

	"github.com/go-pg/pg/orm"
)

const toShortTime = 5 * time.Second

const (
	StatusNew = iota
	StatusSuccess
	StatusFailed
)

type (
	ReminderRepository interface {
		Add(reminder *Reminder) error
		FindActual(timestamp time.Time) ([]Reminder, *pg.Tx, error)
		SetStatus(*Reminder, *pg.Tx) (pg.Result, error)
		Commit(*pg.Tx) error
	}

	Reminder struct {
		ReminderId int       `sql:",pk"`
		CreatedAt  time.Time `pg:"default:now()"`
		Whom       int       `sql:"user_id"`
		What       string    `sql:"content",pg:"notnull"`
		NotifyAt   time.Time
		When       time.Duration `sql:"-"`
		Status     int
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

func (r *Reminder) IsItToday() bool {
	return time.Now().Add(r.When).Day() == time.Now().Day()
}
