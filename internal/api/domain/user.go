package domain

type (
	UserResponse struct {
		Response []User
	}

	User struct {
		Id              int32
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		IsClosed        bool   `json:"is_closed"`
		CanAccessClosed bool   `json:"can_access_closed"`
		TimeZone        int32  `json:"timezone"`
	}
)
