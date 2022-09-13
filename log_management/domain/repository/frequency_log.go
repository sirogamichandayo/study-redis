package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"time"
)

type FrequencyLog struct {
}

func (fl *FrequencyLog) GetUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level domain.LogLevel) (*time.Time, error) {
	t, err := cmd.Get(ctx, fl.updatedAtKey(name, level)).Time()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (fl *FrequencyLog) SetUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level domain.LogLevel, u *domain.FrequencyLogUpdatedAt) error {
	key := fl.updatedAtKey(name, level)
	return cmd.Set(ctx, key, u.String(), 0).Err()
}

func (fl *FrequencyLog) IncrCount(ctx context.Context, cmd redis.Cmdable, lm *domain.LogMessage) error {
	return cmd.ZIncrBy(ctx, fl.countKey(lm.Name(), lm.Level()), 1, lm.Message()).Err()
}

func (fl *FrequencyLog) Archive(ctx context.Context, cmd redis.Cmdable, name string, level domain.LogLevel) (bool, error) {
	ok, err := cmd.RenameNX(ctx,
		fl.updatedAtKey(name, level),
		fl.archiveUpdatedAtKey(name, level),
	).Result()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return cmd.RenameNX(ctx,
		fl.countKey(name, level),
		fl.archiveCountKey(name, level),
	).Result()
}

func (fl *FrequencyLog) WatchMakeAtKey(ctx context.Context, client redis.UniversalClient, fn func(*redis.Tx) error, name string, level domain.LogLevel) error {
	return client.Watch(ctx, fn, fl.updatedAtKey(name, level))
}

func (*FrequencyLog) updatedAtKey(name string, level domain.LogLevel) string {
	return fmt.Sprintf(frequencyLogUpdatedAtFormat, name, level.String())
}

func (*FrequencyLog) countKey(name string, level domain.LogLevel) string {
	return fmt.Sprintf(frequencyLogCountFormat, name, level.String())
}

func (*FrequencyLog) archiveUpdatedAtKey(name string, level domain.LogLevel) string {
	return fmt.Sprintf(frequencyLogArchiveUpdatedAtFormat, name, level.String())
}

func (*FrequencyLog) archiveCountKey(name string, level domain.LogLevel) string {
	return fmt.Sprintf(frequencyLogArchiveCountFormat, name, level.String())
}

const (
	frequencyLogUpdatedAtFormat        = "frequency-log:%s:%s:updated-at"
	frequencyLogCountFormat            = "frequency-log:%s:%s:count"
	frequencyLogArchiveUpdatedAtFormat = "frequency-log:%s:%s:archive:updated-at"
	frequencyLogArchiveCountFormat     = "frequency-log:%s:%s:archive:count"
)
