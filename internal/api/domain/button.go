package domain

import (
	"encoding/json"

	"github.com/sepuka/campaner/internal/command/domain"
)

type (
	ButtonPayload struct {
		Button  string `json:"button"`
		Command string `json:"command"`
	}
)

func (obj ButtonPayload) String() string {
	var str []byte
	str, _ = json.Marshal(obj)

	return string(str)
}

func (obj ButtonPayload) IsChangeTaskButton() bool {
	for _, commandId := range domain.ShiftedTasks {
		if commandId.String() == obj.Command {
			return true
		}
	}

	return false
}

func (obj ButtonPayload) IsStartButton() bool {
	return obj.Command == domain.CommandStartId.String()
}

func (obj ButtonPayload) IsOKButton() bool {
	return obj.Command == domain.CommandOkId.String()
}
