package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
)

type FrequencyLog struct {
}

func (fl *FrequencyLog) GetMakeAt(ctx context.Context, cmd redis.Cmdable, name string, level domain.LogLevel) (string, error) {
	str, err := cmd.Get(ctx, fl.CommonKey(name, level)).Result()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (fl *FrequencyLog) SetMakeAt(ctx context.Context, cmd redis.Cmdable, lm *domain.LogMessage) error {
	key := fl.CommonStartKey(lm.Name(), lm.Level())
	return cmd.Set(ctx, key, lm.MakeAt().String(), 0).Err()
}

func (fl *FrequencyLog) IncrFrequencyCount(ctx context.Context, cmd redis.Cmdable, lm *domain.LogMessage) error {
	return cmd.ZIncrBy(ctx, fl.CommonKey(lm.Name(), lm.Level()), 1, lm.Message()).Err()
}

func (fl *FrequencyLog) Archive(ctx context.Context, cmd redis.Cmdable, name string, level domain.LogLevel) (bool, error) {
	key := fl.CommonKey(name, level)
	ok, err := cmd.RenameNX(ctx, key, key+":last").Result()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	startKey := fl.CommonStartKey(name, level)
	return cmd.RenameNX(ctx, startKey, key+":pstart").Result()
}

func (fl *FrequencyLog) RenameToLastKey(ctx context.Context, cmd redis.Cmdable, name string, level domain.LogLevel) (bool, error) {
	key := fl.CommonKey(name, level)
	return cmd.RenameNX(ctx, key, key+":last").Result()
}

func (*FrequencyLog) CommonKey(name string, level domain.LogLevel) string {
	return fmt.Sprintf("common:%s:%s", name, level.String())
}

func (*FrequencyLog) CommonStartKey(name string, level domain.LogLevel) string {
	return fmt.Sprintf("common:%s:%s:start", name, level.String())
}
