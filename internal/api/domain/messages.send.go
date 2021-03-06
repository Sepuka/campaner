package domain

type (
	ButtonText string
	ButtonType string
	Action     struct {
		Type    ButtonType `json:"type"`
		Label   ButtonText `json:"label"`
		Payload string     `json:"payload"`
	}

	Button struct {
		Action Action `json:"action"`
		Color  string `json:"color"`
	}

	Keyboard struct {
		OneTime bool       `json:"one_time"`
		Buttons [][]Button `json:"buttons"`
	}

	MessagesSend struct {
		Keyboard    string `url:"keyboard"`
		Message     string `url:"message"`
		AccessToken string `url:"access_token"`
		ApiVersion  string `url:"v"`
		PeerId      int    `url:"peer_id"`
		RandomId    int64  `url:"random_id"`
	}
)
