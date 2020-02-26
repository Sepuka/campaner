package analyzer

import "time"

const toShortTime = 5 * time.Second

type Reminder struct {
	when time.Duration
	whom int
	what string
}

func NewReminder(whom int) *Reminder {
	return &Reminder{
		whom: whom,
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

func (r *Reminder) When() time.Duration {
	return r.when
}

func (r *Reminder) What() string {
	return r.what
}

func (r *Reminder) IsImmediate() bool {
	return r.when < toShortTime
}
