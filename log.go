package uptimerobot

type LogType int

const (
	LogTypeDown    LogType = 1
	LogTypeUp              = 2
	LogTypePaused          = 99
	LogTypeStarted         = 98
)

type Log struct {
	Type     LogType         `json:"type,string"`
	DateTime UptimeRobotDate `json:"datetime"`
}
