package payload

import (
	"strconv"
	"time"

	domain2 "github.com/sepuka/campaner/internal/command/domain"

	domainApi "github.com/sepuka/campaner/internal/api/domain"

	"github.com/sepuka/campaner/internal/domain"
)

func HandleChangeTaskPayload(msg string, payload domainApi.ButtonPayload, reminder *domain.Reminder) error {
	var (
		taskId int64
		err    error
	)

	if taskId, err = strconv.ParseInt(payload.Button, 10, 64); err != nil {
		return err
	}

	switch payload.Command {
	case domain2.CommandCancelId.String():
		reminder.ReminderId = int(taskId)
		reminder.Status = domain.StatusCanceled
	case domain2.CommandAfter15MinutesId.String():
		reminder.Status = domain.StatusCopied
		reminder.ReminderId = int(taskId)
		reminder.When = time.Duration(15) * time.Minute
	case domain2.CommandOnTheEveId.String(), domain2.CommandBefore5Minutes.String():
		reminder.ReminderId = int(taskId)
		reminder.Status = domain.StatusShifted
		reminder.RewriteSubject(payload.Command) //TODO
	}

	return nil
}
