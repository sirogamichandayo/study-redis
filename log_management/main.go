package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log_management/adapter/kvs"
	"log_management/domain"
	"log_management/domain/repository"
	"time"
)

func main() {
	ctx := context.Background()
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := c.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	fl := kvs.NewFrequencyLogImpl()
	lm := domain.NewLogMessage(
		"name",
		"this is test",
		domain.Warning,
	)

	err := StoreLog(ctx, fl, c, lm)
	if err != nil {
		panic(err)
	}
}

func StoreLog(ctx context.Context, fl repository.FrequencyLogInterface, client *redis.Client, lm *domain.LogMessage) error {
	name, level := lm.Name(), lm.Level()

	txf := func(tx *redis.Tx) error {
		rawUpdatedAt, sErr := fl.GetUpdatedAt(ctx, tx, name, level)
		if sErr != nil && sErr != redis.Nil {
			return sErr
		}
		if sErr == redis.Nil {
			rawUpdatedAt = time.Now()
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
				uOk, err := fl.ArchiveUpdatedAt(ctx, pipe, name, level)
				if err != nil {
					return err
				}
				if !uOk {
					return errors.New("failed archive updated at")
				}
				cOk, err := fl.ArchiveCount(ctx, pipe, name, level)
				if err != nil {
					return err
				}
				if !cOk {
					return errors.New("failed archive count")
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

	for i := 0; i < 100; i++ {
		err := fl.WatchMakeAtKey(ctx, client, txf, name, level)
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
