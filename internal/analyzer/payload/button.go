package payload

import (
	"encoding/json"
	"strconv"

	"github.com/sepuka/campaner/internal/api/domain"
)

func GetTaskId(raw string) (taskId int64, err error) {
	var (
		payload domain.ButtonPayload
	)

	if err = json.Unmarshal([]byte(raw), &payload); err != nil {
		return 0, err
	}

	taskId, err = strconv.ParseInt(payload.Button, 10, 64)

	return taskId, err
}
