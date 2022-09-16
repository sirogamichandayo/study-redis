package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"log_management/domain/repository"
	"log_management/domain/repository/model"
)

func GetAllLogs(ctx context.Context, fl repository.FrequencyLogInterface, client *redis.Client, name string, level *domain.LogLevel) ([]*domain.FrequencyLogCount, error) {
	op := model.NewAllRankOption()

	countModels, err := fl.GetCountsByRank(ctx, client, name, level, op)
	if err != nil {
		return nil, err
	}

	return domain.MakeFrequencyLogFromModelList(name, level, countModels), nil
}
