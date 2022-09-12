package domain

import "time"

const makeAtFormat = "2006-01-02"

type LogMessageMakeAt struct {
	t time.Time
}

func (mt *LogMessageMakeAt) Time() time.Time {
	return mt.t
}

func (mt *LogMessageMakeAt) String() string {
	return mt.t.Format(makeAtFormat)
}

func (mt *LogMessageMakeAt) Before(t *LogMessageMakeAt) bool {
	return mt.t.Before(t.t)
}

func ParseLogMessageMakeAt(value string) (*LogMessageMakeAt, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation(makeAtFormat, value, jst)
	return &LogMessageMakeAt{t}, nil
}

type LogMessage struct {
	name    string
	message string
	level   LogLevel
	makeAt  *LogMessageMakeAt
}

func NewLogMessage(name string, message string, severity LogLevel) *LogMessage {
	n := time.Now()
	return &LogMessage{
		name, message, severity,
		&LogMessageMakeAt{time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())},
	}
}

func Reconstruct(name string, message string, level LogLevel, makeAt string) (*LogMessage, error) {
	m, err := ParseLogMessageMakeAt(makeAt)
	if err != nil {
		return nil, err
	}
	return &LogMessage{
		name, message, level, m,
	}, nil
}

func (l LogMessage) Name() string {
	return l.name
}

func (l LogMessage) Message() string {
	return l.message
}

func (l LogMessage) Level() LogLevel {
	return l.level
}

func (l LogMessage) MakeAt() *LogMessageMakeAt {
	return l.makeAt
}

// Before is l.makeAt > mt
func (l LogMessage) Before(mt *LogMessageMakeAt) bool {
	return l.makeAt.Before(mt)
}
