package domain

type (
	CommandId      string
	taskCommandIds []CommandId
)

var (
	ShiftedTasks = taskCommandIds{
		CommandAfter15MinutesId,
		CommandOnTheEveId,
		CommandBefore5Minutes,
		CommandCancelId,
	}
)

const (
	CommandStartId          CommandId = `start`
	CommandOkId             CommandId = `OK`
	CommandAfter15MinutesId CommandId = `after15minutes`
	CommandOnTheEveId       CommandId = `onTheEve`
	CommandBefore5Minutes   CommandId = `before5Minutes`
	CommandCancelId         CommandId = `cancel`
)

func (t CommandId) String() string {
	return string(t)
}
