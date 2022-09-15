package domain

import (
	"time"
)

const makeAtFormat = time.RFC3339

type LogMessageMakeAt struct {
	t time.Time
}

func (mt *LogMessageMakeAt) Time() time.Time {
	return mt.t
}

func (mt *LogMessageMakeAt) String() string {
	return mt.t.Format(makeAtFormat)
}

type FrequencyLogUpdatedAt struct {
	t time.Time
}

func NewFrequencyLogUpdatedAt(t time.Time) (*FrequencyLogUpdatedAt, error) {
	ti := time.Date(
		t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60),
	)
	return &FrequencyLogUpdatedAt{ti}, nil
}

func (mt *FrequencyLogUpdatedAt) Time() time.Time {
	return mt.t
}

func (mt *FrequencyLogUpdatedAt) ShouldArchive(now time.Time) (bool, error) {
	t := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	return mt.t.Before(t), nil
}

type LogMessage struct {
	name    string
	message string
	level   *LogLevel
	makeAt  *LogMessageMakeAt
}

func NewLogMessage(name string, message string, severity LogLevel) *LogMessage {
	n := time.Now()
	t := time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	return &LogMessage{
		name, message, &severity,
		&LogMessageMakeAt{t},
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
