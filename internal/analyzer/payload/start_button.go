package payload

import (
	"github.com/sepuka/campaner/internal/api/domain"
	domain2 "github.com/sepuka/campaner/internal/domain"
)

func HandleStartButtonPayload(payload domain.ButtonPayload, reminder *domain2.Reminder) error {
	reminder.Status = domain2.StatusBarren
	reminder.RewriteSubject(`напишите фразу, а мы попробуем разобраться когда и о чём вам нужно напомнить`)

	return nil
}
