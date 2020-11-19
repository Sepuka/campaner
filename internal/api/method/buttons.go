package method

import (
	"strconv"

	"github.com/sepuka/campaner/internal/api/domain"
)

const (
	TextButtonType   domain.ButtonType = `text`
	CancelButton     domain.ButtonText = `cancel`
	OKButton         domain.ButtonText = `OK`
	Later15MinButton domain.ButtonText = `через 15 минут`
	Before5Minutes   domain.ButtonText = `за 5 минут`
	OnTheEve         domain.ButtonText = `накануне вечером`
)

func cancel(cancelId int) [][]domain.Button {
	return [][]domain.Button{
		{
			{
				Color: `negative`,
				Action: domain.Action{
					Type:    TextButtonType,
					Label:   CancelButton,
					Payload: domain.ButtonPayload{Button: strconv.Itoa(cancelId)}.String(),
				},
			},
			{
				Color: `secondary`,
				Action: domain.Action{
					Type:    TextButtonType,
					Label:   Before5Minutes,
					Payload: domain.ButtonPayload{Button: strconv.Itoa(cancelId)}.String(),
				},
			},
		},
	}
}

func cancelWithEve(cancelId int) [][]domain.Button {
	return [][]domain.Button{
		{
			{
				Color: `negative`,
				Action: domain.Action{
					Type:    TextButtonType,
					Label:   CancelButton,
					Payload: domain.ButtonPayload{Button: strconv.Itoa(cancelId)}.String(),
				},
			},
			{
				Color: `secondary`,
				Action: domain.Action{
					Type:    TextButtonType,
					Label:   OnTheEve,
					Payload: domain.ButtonPayload{Button: strconv.Itoa(cancelId)}.String(),
				},
			},
		},
	}
}

func delayAndOk(remindId int) [][]domain.Button {
	return [][]domain.Button{
		{
			{
				Color: `positive`,
				Action: domain.Action{
					Type:    TextButtonType,
					Label:   Later15MinButton,
					Payload: domain.ButtonPayload{Button: strconv.Itoa(remindId)}.String(),
				},
			},
			{
				Color: `primary`,
				Action: domain.Action{
					Type:    TextButtonType,
					Label:   OKButton,
					Payload: domain.ButtonPayload{Button: strconv.Itoa(remindId)}.String(),
				},
			},
		},
	}
}
