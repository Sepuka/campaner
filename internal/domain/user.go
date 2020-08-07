package domain

import "time"

type (
	UserRepository interface {
		Get(userId int) (*User, error)
	}

	User struct {
		Id        int `sql:"user_id,pk"`
		Timezone  int
		UpdatedAt time.Time
	}
)
