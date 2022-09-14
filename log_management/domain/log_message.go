package domain

import "time"

const makeAtFormat = time.RFC3339

type LogMessageMakeAt struct {
	t *time.Time
}

func (mt *LogMessageMakeAt) Time() *time.Time {
	return mt.t
}

func (mt *LogMessageMakeAt) String() string {
	return mt.t.Format(makeAtFormat)
}

const frequencyLogUpdatedAtFormat = "2006-01-02 15"

type FrequencyLogUpdatedAt struct {
	t time.Time
}

func NewFrequencyLogUpdatedAt(t *time.Time) (*FrequencyLogUpdatedAt, error) {
	jst, e := time.LoadLocation("Asia/Tokyo")
	if e != nil {
		return nil, e
	}
	ti := time.Date(
		t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, jst,
	)
	return &FrequencyLogUpdatedAt{ti}, nil
}

func (mt *FrequencyLogUpdatedAt) String() string {
	return mt.t.Format(frequencyLogUpdatedAtFormat)
}

func (mt *FrequencyLogUpdatedAt) ShouldArchive() bool {
	return true
	// return mt.t.After(at.t)
}

type LogMessage struct {
	name    string
	message string
	level   *LogLevel
	makeAt  *LogMessageMakeAt
}

func NewLogMessage(name string, message string, severity LogLevel) *LogMessage {
	n := time.Now()
	t := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	return &LogMessage{
		name, message, &severity,
		&LogMessageMakeAt{&t},
	}
}

func (l LogMessage) Name() string {
	return l.name
}

func (l LogMessage) Message() string {
	return l.message
}

func (l LogMessage) Level() *LogLevel {
	return l.level
}

func (l LogMessage) MakeAt() *LogMessageMakeAt {
	return l.makeAt
}
