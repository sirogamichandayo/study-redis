package domain

type LogLevel string

const (
	Debug    LogLevel = "debug"
	Info     LogLevel = "info"
	Warning  LogLevel = "warning"
	Error    LogLevel = "error"
	Critical LogLevel = "critical"
)

func (l LogLevel) String() string {
	return string(l)
}
