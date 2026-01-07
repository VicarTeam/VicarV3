package auth

import (
	"context"
	"fmt"
	"time"
	"vicar-backend/cache"
	"vicar-backend/log"
)

func registerIdentityTokenInCache(userID, deviceName, token string, exp int64) error {
	return cache.Redis.Set(context.Background(), fmt.Sprintf("jwt-identity:%s:%s:%s", userID, deviceName, token), "true", time.Hour*24*30).Err()
}

func unregisterAllIdentityTokensForDeviceInCache(userID, deviceName string) error {
	keys := cache.Redis.Keys(context.Background(), fmt.Sprintf("jwt-identity:%s:%s:*", userID, deviceName)).Val()
	if len(keys) == 0 {
		return nil
	}

	return cache.Redis.Del(context.Background(), keys...).Err()
}

func unregisterAllIdentityTokensInCache(userID string) error {
	keys := cache.Redis.Keys(context.Background(), fmt.Sprintf("jwt-identity:%s:*", userID)).Val()
	if len(keys) == 0 {
		return nil
	}

	return cache.Redis.Del(context.Background(), keys...).Err()
}

func isIdentityTokenValidInCache(userID, deviceName, token string) (bool, error) {
	log.Debug(log.Auth, "🔑", "Checking if identity token is valid in cache for user %s, device %s, token %s", userID, deviceName, token)
	i, err := cache.Redis.Exists(context.Background(), fmt.Sprintf("jwt-identity:%s:%s:%s", userID, deviceName, token)).Result()
	if err != nil {
		return false, err
	}

	return i > 0, nil
}
