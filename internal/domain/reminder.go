package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sepuka/campaner/internal/speeches"

	"github.com/go-pg/pg"

	"github.com/go-pg/pg/orm"
)

const toShortTime = 5 * time.Second

const (
	StatusNew = iota
	StatusSuccess
	StatusFailed
	StatusCanceled
	StatusCopied
)

type (
	ReminderRepository interface {
		FindActual(timestamp time.Time) ([]Reminder, *pg.Tx, error)
		SetStatus(*Reminder, *pg.Tx) (pg.Result, error)
		Commit(*pg.Tx) error
		Scheduled(userId int, limit uint32) ([]Reminder, error)
	}

	TaskManager interface {
		Add(reminder *Reminder) error
		Cancel(reminder *Reminder) error
		Copy(reminder *Reminder) error
	}

	Reminder struct {
		ReminderId int       `sql:",pk"`
		CreatedAt  time.Time `pg:"default:now()"`
		Whom       int       `sql:"user_id"`
		What       string    `sql:"content",pg:"notnull"`
		NotifyAt   time.Time
		When       time.Duration `sql:"-"`
		Status     int
		Subject    []string `sql:"-"`
	}
)

func NewReminder(whom int) *Reminder {
	return &Reminder{
		Whom: whom,
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

	r.What = r.GetSubject()

	return nil
}

func (r *Reminder) IsImmediate() bool {
	return r.When < toShortTime
}

func (r *Reminder) IsTimeUnknown() bool {
	return r.When == 0
}

func (r *Reminder) RewriteSubject(subject string) {
	r.Subject = []string{subject}
}

func (r *Reminder) AppendSubject(pattern *speeches.Pattern) {
	r.Subject = append(r.Subject, pattern.String())
}

func (r *Reminder) GetSubject() string {
	return strings.Join(r.Subject, ` `)
}

func (r *Reminder) String() string {
	return fmt.Sprintf(`"%s" at %s`, r.What, r.NotifyAt.Format(`2006-01-02 15:04:05`))
}
