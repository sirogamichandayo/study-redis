package test_tools

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

func MakeFakeClient() (*redis.Client, error) {
	r, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr(),
		Password: "",
		DB:       0,
	})
	return client, nil
}
