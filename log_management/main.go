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

func log_common(ctx context.Context, client *redis.Client, lm *domain.LogMessage) error {
	name, level := lm.Name(), lm.Level()
	fl := repository.FrequencyLog{}

	txf := func(tx *redis.Tx) error {
		str, getErr := fl.GetMakeAt(ctx, tx, name, level)
		if getErr != nil {
			return getErr
		}

		if getErr != redis.Nil {
			storedMakeAt, parseErr := domain.ParseLogMessageMakeAt(str)
			if parseErr != nil {
				return parseErr
			}

			if lm.Before(storedMakeAt) {
				ok, aErr := fl.Archive(ctx, tx, name, level)
				if aErr != nil {
					return aErr
				}
				if !ok {
					return errors.New("failed archiving")
				}
				if sErr := fl.SetMakeAt(ctx, tx, lm); sErr != nil {
					return sErr
				}
			}
		}

		return fl.IncrFrequencyCount(ctx, tx, lm)
	}

	for i := 0; i < 100; i++ {
		err := client.Watch(ctx, txf, fl.CommonStartKey(name, level))
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
