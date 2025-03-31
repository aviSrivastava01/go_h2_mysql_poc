package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	myredis "go-redis-mysql-poc/redis" // Replace with your actual import path

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCmdable is a mock implementation of the redis.Cmdable interface.
type MockCmdable struct {
	mock.Mock
}

func (m *MockCmdable) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockCmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockCmdable) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockCmdable) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

func TestGet(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"
	expectedValue := "testvalue"

	stringCmd := redis.NewStringCmd(context.Background())
	stringCmd.SetVal(expectedValue)

	mockCmdable.On("Get", myredis.ctx, key).Return(stringCmd)

	value, err := myredis.Get(key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)

	mockCmdable.AssertExpectations(t)
}

func TestGet_Nil(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"

	stringCmd := redis.NewStringCmd(context.Background())
	stringCmd.SetErr(redis.Nil)

	mockCmdable.On("Get", myredis.ctx, key).Return(stringCmd)

	value, err := myredis.Get(key)

	assert.NoError(t, err)
	assert.Equal(t, "", value)

	mockCmdable.AssertExpectations(t)
}

func TestGet_Error(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"
	expectedError := errors.New("redis error")

	stringCmd := redis.NewStringCmd(context.Background())
	stringCmd.SetErr(expectedError)

	mockCmdable.On("Get", myredis.ctx, key).Return(stringCmd)

	value, err := myredis.Get(key)

	assert.Error(t, err)
	assert.Equal(t, "", value)
	assert.Equal(t, expectedError, err)

	mockCmdable.AssertExpectations(t)
}

func TestSet(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"
	value := "testvalue"
	expiration := 10

	statusCmd := redis.NewStatusCmd(context.Background())
	statusCmd.SetVal("OK")

	mockCmdable.On("Set", myredis.ctx, key, value, time.Duration(expiration)*time.Second).Return(statusCmd)

	err := myredis.Set(key, value, expiration)

	assert.NoError(t, err)

	mockCmdable.AssertExpectations(t)
}

func TestSet_Error(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"
	value := "testvalue"
	expiration := 10
	expectedError := errors.New("redis error")

	statusCmd := redis.NewStatusCmd(context.Background())
	statusCmd.SetErr(expectedError)

	mockCmdable.On("Set", myredis.ctx, key, value, time.Duration(expiration)*time.Second).Return(statusCmd)

	err := myredis.Set(key, value, expiration)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockCmdable.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"

	intCmd := redis.NewIntCmd(context.Background())
	intCmd.SetVal(1)

	mockCmdable.On("Del", myredis.ctx, []string{key}).Return(intCmd)

	err := myredis.Delete(key)

	assert.NoError(t, err)

	mockCmdable.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	mockCmdable := new(MockCmdable)

	// Save and restore the original RedisClient
	originalRedisClient := myredis.RedisClient
	defer func() {
		myredis.RedisClient = originalRedisClient
	}()

	myredis.RedisClient = mockCmdable // Assign the mock directly

	key := "testkey"
	expectedError := errors.New("redis error")

	intCmd := redis.NewIntCmd(context.Background())
	intCmd.SetErr(expectedError)

	mockCmdable.On("Del", myredis.ctx, []string{key}).Return(intCmd)

	err := myredis.Delete(key)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockCmdable.AssertExpectations(t)
}
