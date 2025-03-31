package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient redis.Cmdable // Change the type to redis.Cmdable

var ctx = context.Background()

func ConnectRedis() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	//redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	redisDB := 0
	if redisDBStr != "" {
		var err error
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			log.Printf("Error converting REDIS_DB to integer, using default DB 0: %v", err)
			redisDB = 0 // Default to database 0 if conversion fails
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",      // no password set
		DB:       redisDB, // use default DB
	})

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis!")
}

func Get(key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist (not an error)
	} else if err != nil {
		return "", err // Actual error
	}
	return val, nil
}

func Set(key string, value interface{}, expiration int) error {
	err := RedisClient.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func Delete(key string) error {
	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
