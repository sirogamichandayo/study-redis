//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE
package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"log_management/domain/repository/model"
	"time"
)

type FrequencyLogInterface interface {
	GetUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel) (time.Time, error)
	SetUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel, u *domain.FrequencyLogUpdatedAt) error
	IncrCount(ctx context.Context, cmd redis.Cmdable, lm *domain.LogMessage) error
	GetCountsByRank(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel, op *model.RankOption) ([]*model.LogMessageCount, error)
	ArchiveUpdatedAt(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel) error
	ArchiveCount(ctx context.Context, cmd redis.Cmdable, name string, level *domain.LogLevel) error
	WatchUpdatedAt(ctx context.Context, client redis.UniversalClient, fn func(*redis.Tx) error, name string, level *domain.LogLevel) error
}
