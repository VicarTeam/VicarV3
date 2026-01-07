package cache

import (
	"context"
	"os"
	"time"
	"vicar-backend/log"
	"vicar-backend/util"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
)

func Initialize() error {
	log.Info(log.Cache, "🔌", "Initializing cache...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	log.Success(log.Cache, "✅", "Cache initialized successfully")

	Redis = rdb

	initializeLocking()

	return nil
}

// Creates a cached state for the given data.
func BeginState(category string, payload any, duration time.Duration) string {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	stateKey, err := util.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}

	Redis.Set(context.Background(), "state:"+category+":"+stateKey, string(data), duration)

	return stateKey
}

func BeginSpecificState(category string, state string, payload any, duration time.Duration) {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	Redis.Set(context.Background(), "state:"+category+":"+state, string(data), duration)
}

// Deletes a cached state.
func EndState(category string, state string, payload any) bool {
	data, err := Redis.Get(context.Background(), "state:"+category+":"+state).Result()
	if err != nil {
		if err == redis.Nil {
			return false
		}

		panic(err)
	}

	err = json.Unmarshal([]byte(data), payload)
	if err != nil {
		panic(err)
	}

	Redis.Del(context.Background(), "state:"+category+":"+state)

	return true
}

// Retrieves a cached state. This will not delete the state. Use EndState to delete/consume the state.
func GetState(category string, state string, payload any) bool {
	data, err := Redis.Get(context.Background(), "state:"+category+":"+state).Result()
	if err != nil {
		if err == redis.Nil {
			return false
		}

		panic(err)
	}

	err = json.Unmarshal([]byte(data), payload)
	if err != nil {
		panic(err)
	}

	return true
}

// Checks if a cached state exists.
func HasState(category string, state string) bool {
	_, err := Redis.Get(context.Background(), "state:"+category+":"+state).Result()
	if err != nil {
		if err == redis.Nil {
			return false
		}

		panic(err)
	}

	return true
}

func GetJSON(key string, payload any) error {
	data, err := Redis.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), payload)
}
