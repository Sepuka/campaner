package payload

import (
	domainApi "github.com/sepuka/campaner/internal/api/domain"
	"github.com/sepuka/campaner/internal/domain"
)

func HandleOKButtonPayload(msg string, payload domainApi.ButtonPayload, reminder *domain.Reminder) error {
	reminder.Status = domain.StatusBarren

	return nil
}
