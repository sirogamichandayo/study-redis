package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"log_management/domain/repository"
	redTime "log_management/tools/red_time"
)

const storeLogMaxRetries = 30

func StoreLog(ctx context.Context, fl repository.FrequencyLogInterface, client *redis.Client, lm *domain.LogMessage, redTime redTime.ITime) error {
	name, level := lm.Name(), lm.Level()

	for i := 0; i < storeLogMaxRetries; i++ {
		err := fl.WatchUpdatedAt(
			ctx,
			client,
			makeStoreLogFunc(ctx, fl, lm, redTime),
			name,
			level,
		)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}

	return errors.New("increment reached maximum number of retries")
}

func makeStoreLogFunc(
	ctx context.Context,
	fl repository.FrequencyLogInterface,
	lm *domain.LogMessage,
	redTime redTime.ITime,
) func(tx *redis.Tx) error {
	name, level := lm.Name(), lm.Level()

	return func(tx *redis.Tx) error {
		rawUpdatedAt, sErr := fl.GetUpdatedAt(ctx, tx, name, level)
		if sErr != nil && sErr != redis.Nil {
			return sErr
		}
		if sErr == redis.Nil {
			rawUpdatedAt = redTime.Now()
		}
		updatedAt, uErr := domain.NewFrequencyLogUpdatedAt(rawUpdatedAt)
		if uErr != nil {
			return nil
		}

		_, pErr := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if sErr == redis.Nil {
				if suErr := fl.SetUpdatedAt(ctx, pipe, name, level, updatedAt); suErr != nil {
					return suErr
				}
			}

			shouldArchive, sErr := updatedAt.ShouldArchive(lm.MakeAt().Time())
			if sErr != nil {
				return sErr
			}
			if shouldArchive {
				if err := fl.ArchiveUpdatedAt(ctx, pipe, name, level); err != nil {
					return err
				}
				if err := fl.ArchiveCount(ctx, pipe, name, level); err != nil {
					return err
				}
				newUpdatedAt, aErr := domain.NewFrequencyLogUpdatedAt(lm.MakeAt().Time())
				if aErr != nil {
					return aErr
				}
				if sErr := fl.SetUpdatedAt(ctx, pipe, name, level, newUpdatedAt); sErr != nil {
					return sErr
				}
			}

			return fl.IncrCount(ctx, pipe, lm)
		})
		return pErr
	}
}
