package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"time"
)

type FrequencyLogInterface interface {
	GetUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel) (*time.Time, error)
	SetUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel, u *domain.FrequencyLogUpdatedAt) error
	IncrCount(ctx context.Context, cmd redis.Cmdable, lm *domain.LogMessage) error
	ArchiveUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel) (bool, error)
	ArchiveCount(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel) (bool, error)
	WatchMakeAtKey(ctx context.Context, client redis.UniversalClient, fn func(*redis.Tx) error, name string, level *domain.LogLevel) error
}
