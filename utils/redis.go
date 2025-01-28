package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	configs, _ := LoadConfig("./../")
	fmt.Println(configs.RedisHost, "333")
	client := redis.NewClient(&redis.Options{
		Addr:     configs.RedisHost,
	})
	fmt.Println(client, "xxc")
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value) // 🔍 Convert value to JSON
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	cachedData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	// 🔍 Unmarshal into destination struct
	return json.Unmarshal([]byte(cachedData), dest)
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}

// InvalidateAllCampaignCaches removes all campaign-related cache keys
func (r *RedisCache) InvalidateAllCampaignCaches() {
	ctx := context.Background()

	// Define patterns for campaign cache keys
	keyPatterns := []string{
		"campaigns_*",  // Matches all campaign keys
	}

	for _, pattern := range keyPatterns {
		iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			err := r.client.Del(ctx, iter.Val()).Err()
			if err != nil {
				fmt.Println("Error deleting cache key:", iter.Val(), err)
			} else {
				fmt.Println("Deleted cache key:", iter.Val())
			}
		}

		if err := iter.Err(); err != nil {
			fmt.Println("Error scanning Redis keys:", err)
		}
	}
}