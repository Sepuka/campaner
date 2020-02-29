package domain

import "time"

const toShortTime = 5 * time.Second

type Reminder struct {
	when time.Duration
	whom int
	what string
}

func NewReminder(whom int, what string, when time.Duration) *Reminder {
	return &Reminder{
		whom: whom,
		what: what,
		when: when,
	}
}

func NewImmediateReminder(
	whom int,
	what string,
) *Reminder {
	return &Reminder{
		when: time.Nanosecond,
		whom: whom,
		what: what,
	}
}

func (r *Reminder) Whom() int {
	return r.whom
}

func (r *Reminder) SetWhen(when time.Duration) {
	r.when = when
}

func (r *Reminder) When() time.Duration {
	return r.when
}

func (r *Reminder) What() string {
	return r.what
}

func (r *Reminder) IsImmediate() bool {
	return r.when < toShortTime
}

func (r *Reminder) IsValid() bool {
	return r.when > 0
}
