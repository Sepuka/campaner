package domain

import "encoding/json"

type (
	ButtonPayload struct {
		Button string `json:"button"`
	}
)

func (obj ButtonPayload) String() string {
	var str []byte
	str, _ = json.Marshal(obj)

	return string(str)
}
