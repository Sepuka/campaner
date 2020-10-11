package repository

import (
	"time"

	"github.com/sepuka/campaner/internal/errors"

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

func (r *ReminderRepository) Cancel(reminder *domain.Reminder) error {
	var (
		err error
	)

	err = r.
		db.
		Model(reminder).
		Where(`reminder_id = ? AND user_id = ? AND status = ?`, reminder.ReminderId, reminder.Whom, domain.StatusNew).
		Select()

	if err != nil {
		return err
	}

	reminder.Status = domain.StatusCanceled

	return r.
		db.
		Update(reminder)
}

func (r *ReminderRepository) Copy(reminder *domain.Reminder) error {
	var (
		donor = &domain.Reminder{
			ReminderId: reminder.ReminderId,
		}
		err error
	)

	err = r.
		db.
		Model(donor).
		WherePK().
		Select()
	if err != nil {
		return err
	}

	if donor.Whom != reminder.Whom {
		return errors.NewTaskManagerError(`wrong user id`)
	}

	if donor.Status != domain.StatusSuccess {
		return errors.NewTaskManagerError(`wrong status`)
	}

	reminder.RewriteSubject(donor.What)
	reminder.ReminderId = 0
	reminder.Status = domain.StatusNew

	return r.Add(reminder)
}
