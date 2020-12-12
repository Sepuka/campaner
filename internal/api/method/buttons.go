package method

import (
	"strconv"

	domain2 "github.com/sepuka/campaner/internal/command/domain"

	"github.com/sepuka/campaner/internal/api/domain"
)

const (
	TextButtonType   domain.ButtonType = `text`
	CancelButton     domain.ButtonText = `отменить`
	OKButton         domain.ButtonText = `OK`
	Later15MinButton domain.ButtonText = `через 15 минут`
	Before5Minutes   domain.ButtonText = `за 5 минут`
	OnTheEve         domain.ButtonText = `накануне вечером`
)

func withoutButtons() [][]domain.Button {
	return [][]domain.Button{}
}

func cancel(cancelId int) [][]domain.Button {
	return [][]domain.Button{
		{
			{
				Color: `negative`,
				Action: domain.Action{
					Type:  TextButtonType,
					Label: CancelButton,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(cancelId),
						Command: domain2.CommandCancelId.String(),
					}.String(),
				},
			},
		},
	}
}

func cancelWith5Minutes(cancelId int) [][]domain.Button {
	return [][]domain.Button{
		{
			{
				Color: `negative`,
				Action: domain.Action{
					Type:  TextButtonType,
					Label: CancelButton,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(cancelId),
						Command: domain2.CommandCancelId.String(),
					}.String(),
				},
			},
			{
				Color: `secondary`,
				Action: domain.Action{
					Type:  TextButtonType,
					Label: Before5Minutes,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(cancelId),
						Command: domain2.CommandBefore5Minutes.String(),
					}.String(),
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
					Type:  TextButtonType,
					Label: CancelButton,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(cancelId),
						Command: domain2.CommandCancelId.String(),
					}.String(),
				},
			},
			{
				Color: `secondary`,
				Action: domain.Action{
					Type:  TextButtonType,
					Label: OnTheEve,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(cancelId),
						Command: domain2.CommandOnTheEveId.String(),
					}.String(),
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
					Type:  TextButtonType,
					Label: Later15MinButton,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(remindId),
						Command: domain2.CommandAfter15MinutesId.String(),
					}.String(),
				},
			},
			{
				Color: `primary`,
				Action: domain.Action{
					Type:  TextButtonType,
					Label: OKButton,
					Payload: domain.ButtonPayload{
						Button:  strconv.Itoa(remindId), // TODO may be task ID does not need
						Command: domain2.CommandOkId.String(),
					}.String(),
				},
			},
		},
	}
}
