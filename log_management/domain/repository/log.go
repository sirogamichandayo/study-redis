package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"time"
)

type Log struct {
}

func (l *Log) Push(ctx context.Context, cmd redis.Cmdable, name string, severity domain.LogLevel, message string, t time.Time) error {
	key := l.key(name, severity)
	value := t.Format(time.RFC3339Nano) + " " + message

	return cmd.LPush(ctx, key, value).Err()
}

func (l *Log) Trim(ctx context.Context, cmd redis.Cmdable, name string, severity domain.LogLevel, begin, end int64) error {
	key := l.key(name, severity)
	return cmd.LTrim(ctx, key, begin, end).Err()
}

func (*Log) key(name string, severity domain.LogLevel) string {
	return fmt.Sprintf("recent:%s:%s", name, severity.String())
}
