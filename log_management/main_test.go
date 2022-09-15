package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"log_management/domain"
	mock_repository "log_management/domain/repository/mock"
	"log_management/test_tools"
	mock_redTime "log_management/tools/red_time/mock"
	"testing"
	"time"
)

func Test正常系_StoreLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := "name"
	message := "message"
	level := domain.Warning

	ctx := context.Background()
	flMock := mock_repository.NewMockFrequencyLogInterface(ctrl)
	client, err := test_tools.MakeFakeClient(t)
	assert.Nil(t, err)
	redTimeMock := mock_redTime.NewMockITime(ctrl)
	lm := domain.NewLogMessage(name, message, level)

	tx := func(*redis.Tx) error { return nil }
	gomock.InOrder(
		flMock.EXPECT().WatchMakeAtKey(
			ctx, client,
			gomock.AssignableToTypeOf(tx),
			name, &level).Times(1).Return(nil),
	)

	err = StoreLog(ctx, flMock, client, lm, redTimeMock)
	assert.Nil(t, err)
}

func Test異常系_StoreLog_99回TxFailedErrが帰って来て最後の1回が正常(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := "name"
	message := "message"
	level := domain.Warning

	ctx := context.Background()
	flMock := mock_repository.NewMockFrequencyLogInterface(ctrl)
	client, err := test_tools.MakeFakeClient(t)
	assert.Nil(t, err)
	lm := domain.NewLogMessage(name, message, level)
	redTimeMock := mock_redTime.NewMockITime(ctrl)

	tx := func(*redis.Tx) error { return nil }
	gomock.InOrder(
		flMock.EXPECT().WatchMakeAtKey(ctx, client,
			gomock.AssignableToTypeOf(tx),
			name, &level).Times(99).Return(redis.TxFailedErr),
		flMock.EXPECT().WatchMakeAtKey(ctx, client,
			gomock.AssignableToTypeOf(tx),
			name, &level).Return(nil),
	)

	err = StoreLog(ctx, flMock, client, lm, redTimeMock)
	assert.Nil(t, err)
}

func Test異常系_StoreLog_エラーが返却される(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := "name"
	message := "message"
	level := domain.Warning

	ctx := context.Background()
	flMock := mock_repository.NewMockFrequencyLogInterface(ctrl)
	client, err := test_tools.MakeFakeClient(t)
	assert.Nil(t, err)
	lm := domain.NewLogMessage(name, message, level)
	redTimeMock := mock_redTime.NewMockITime(ctrl)

	tx := func(*redis.Tx) error { return nil }
	gomock.InOrder(
		flMock.EXPECT().WatchMakeAtKey(ctx, client,
			gomock.AssignableToTypeOf(tx),
			name, &level).Times(1).Return(errors.New("test error")),
	)

	err = StoreLog(ctx, flMock, client, lm, redTimeMock)
	assert.EqualError(t, err, "test error")
}

func Test正常系_makeStoreLog初めての呼び出し(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	name := "name"
	message := "message"
	level := domain.Warning
	lm := domain.NewLogMessage(name, message, level)

	flMock := mock_repository.NewMockFrequencyLogInterface(ctrl)
	client, err := test_tools.MakeFakeClient(t)
	assert.Nil(t, err)
	redTimeMock := mock_redTime.NewMockITime(ctrl)

	now := time.Now()
	expectedUpdatedAt, uErr := domain.NewFrequencyLogUpdatedAt(now)
	assert.Nil(t, uErr)

	targetFn := makeStoreLogFunc(ctx, flMock, lm, redTimeMock)

	gomock.InOrder(
		flMock.EXPECT().GetUpdatedAt(ctx, gomock.AssignableToTypeOf(&redis.Tx{}), name, &level).Times(1).
			Return(time.Time{}, redis.Nil),
		redTimeMock.EXPECT().Now().Times(1).
			Return(now),
		flMock.EXPECT().SetUpdatedAt(
			ctx,
			gomock.AssignableToTypeOf(&redis.Pipeline{}),
			name,
			&level,
			test_tools.DiffEq(expectedUpdatedAt, cmp.AllowUnexported(domain.FrequencyLogUpdatedAt{})),
		).Times(1).Return(nil),
		flMock.EXPECT().IncrCount(
			ctx,
			gomock.AssignableToTypeOf(&redis.Pipeline{}),
			lm,
		).Times(1).Return(nil),
	)

	err = client.Watch(ctx, targetFn)
	assert.Nil(t, err)
}

func Test正常系_makeStoreLogアーカイブ処理(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	name := "name"
	message := "message"
	level := domain.Warning
	lm := domain.NewLogMessage(name, message, level)

	flMock := mock_repository.NewMockFrequencyLogInterface(ctrl)
	client, err := test_tools.MakeFakeClient(t)
	assert.Nil(t, err)
	redTimeMock := mock_redTime.NewMockITime(ctrl)

	beforeUpdatedAt := time.Now().Add(-1 * time.Hour)
	newLogUpdatedAt, lErr := domain.NewFrequencyLogUpdatedAt(lm.MakeAt().Time())
	assert.Nil(t, lErr)

	targetFn := makeStoreLogFunc(ctx, flMock, lm, redTimeMock)

	gomock.InOrder(
		flMock.EXPECT().GetUpdatedAt(ctx, gomock.AssignableToTypeOf(&redis.Tx{}), name, &level).
			Times(1).Return(beforeUpdatedAt, nil),
		flMock.EXPECT().ArchiveUpdatedAt(ctx, gomock.AssignableToTypeOf(&redis.Pipeline{}), name, &level).
			Times(1).Return(true, nil),
		flMock.EXPECT().ArchiveCount(ctx, gomock.AssignableToTypeOf(&redis.Pipeline{}), name, &level).
			Times(1).Return(true, nil),
		flMock.EXPECT().SetUpdatedAt(
			ctx,
			gomock.AssignableToTypeOf(&redis.Pipeline{}),
			name,
			&level,
			test_tools.DiffEq(newLogUpdatedAt, cmp.AllowUnexported(domain.FrequencyLogUpdatedAt{})),
		).Times(1).Return(nil),
		flMock.EXPECT().IncrCount(
			ctx,
			gomock.AssignableToTypeOf(&redis.Pipeline{}),
			lm,
		).Times(1).Return(nil),
	)

	err = client.Watch(ctx, targetFn)
	assert.Nil(t, err)
}
