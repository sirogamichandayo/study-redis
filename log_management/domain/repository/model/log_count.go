package model

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

type LogMessageCount struct {
	message string
	count   int
}

func (l LogMessageCount) Message() string {
	return l.message
}

func (l LogMessageCount) Count() int {
	return l.count
}

func MakeLogCountFromZ(z redis.Z) (*LogMessageCount, error) {
	msg, ok := z.Member.(string)
	if !ok {
		return nil, errors.New("failed to cast to string")
	}
	return &LogMessageCount{
		msg,
		int(z.Score),
	}, nil
}

func MakeLogCountListFromZ(zList []redis.Z) ([]*LogMessageCount, error) {
	lcList := make([]*LogMessageCount, len(zList))
	for _, z := range zList {
		lc, mErr := MakeLogCountFromZ(z)
		if mErr != nil {
			return nil, mErr
		}
		lcList = append(lcList, lc)
	}
	return lcList, nil
}
