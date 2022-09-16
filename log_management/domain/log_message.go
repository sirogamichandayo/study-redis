package domain

import (
	"log_management/domain/repository/model"
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

type FrequencyLogCount struct {
	name    string
	message string
	level   *LogLevel
	count   int
}

func (f FrequencyLogCount) Name() string {
	return f.name
}

func (f FrequencyLogCount) Message() string {
	return f.message
}

func (f FrequencyLogCount) Level() *LogLevel {
	return f.level
}

func (f FrequencyLogCount) Count() int {
	return f.count
}

func MakeFrequencyLogFromModel(name string, level *LogLevel, lmc *model.LogMessageCount) *FrequencyLogCount {
	return &FrequencyLogCount{
		name, lmc.Message(), level, lmc.Count(),
	}
}

func MakeFrequencyLogFromModelList(name string, level *LogLevel, lmcList []*model.LogMessageCount) []*FrequencyLogCount {
	lc := make([]*FrequencyLogCount, 0, len(lmcList))
	for _, lmc := range lmcList {
		lc = append(lc, MakeFrequencyLogFromModel(name, level, lmc))
	}
	return lc
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
