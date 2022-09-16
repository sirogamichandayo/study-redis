package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	/*
		if err := c.Ping(ctx).Err(); err != nil {
			panic(err)
		}

		fl := kvs.NewFrequencyLogImpl()
		lm := domain.NewLogMessage(
			"name",
			"this is test",
			domain.Warning,
		)

		err := StoreLog(ctx, fl, c, lm, &redTime.TimeImpl{})
		if err != nil {
			panic(err)
		}

	*/

	res, _ := c.ZRangeWithScores(ctx, "frequency-log:name:warning:count", 0, -1).Result()
	for _, r := range res {
		fmt.Println(r.Member, r.Score)
	}
}
