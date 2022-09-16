package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log_management/domain"
	"log_management/domain/repository"
	"log_management/domain/repository/model"
)

func GetAllLogsByNameAndLevel(ctx context.Context, fl repository.FrequencyLogInterface, client *redis.Client, name string, level *domain.LogLevel) ([]*domain.FrequencyLogCount, error) {
	op := model.NewAllRankOption()

	countModels, err := fl.GetCountsByRank(ctx, client, name, level, op)
	if err != nil {
		return nil, err
	}

	return domain.MakeFrequencyLogFromModelList(name, level, countModels), nil
}

func GetMostFrequentLogByNameAndLevel(ctx context.Context, fl repository.FrequencyLogInterface, client *redis.Client, name string, level *domain.LogLevel) (*domain.FrequencyLogCount, error) {
	op := model.NewMostFrequentOption()

	countModels, err := fl.GetCountsByRank(ctx, client, name, level, op)
	if err != nil {
		return nil, err
	}
	cLen := len(countModels)
	if cLen > 1 {
		return nil, errors.New("program error")
	}
	if cLen == 0 {
		return &domain.FrequencyLogCount{}, nil
	}

	return domain.MakeFrequencyLogFromModel(name, level, countModels[0]), nil
}
