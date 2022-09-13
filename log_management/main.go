package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"log_management/domain/repository"
)

func main() {
}

func StoreLog(ctx context.Context, client *redis.Client, lm *domain.LogMessage) error {
	name, level := lm.Name(), lm.Level()
	fl := repository.FrequencyLog{}

	txf := func(tx *redis.Tx) error {
		rawUpdatedAt, sErr := fl.GetUpdatedAt(ctx, tx, name, level)
		if sErr != nil && sErr != redis.Nil {
			return sErr
		}

		updatedAt, uErr := domain.NewFrequencyLogUpdatedAt(rawUpdatedAt)
		if uErr != nil {
			return nil
		}

		_, pErr := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if updatedAt.ShouldArchive() {
				ok, err := fl.Archive(ctx, pipe, name, level)
				if err != nil {
					return err
				}
				if !ok {
					return errors.New("failed archive")
				}
				newUpdatedAt, aErr := domain.NewFrequencyLogUpdatedAt(lm.MakeAt().Time())
				if aErr != nil {
					return aErr
				}

				return fl.SetUpdatedAt(ctx, pipe, name, level, newUpdatedAt)
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
