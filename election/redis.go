package election

import "github.com/redis/go-redis/v9"

var RedisClient *redis.Client

func InitSimpleRedis(endpoint, pwd string, db int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     endpoint,
		Password: pwd,
		DB:       db,
	})
}

func InitRedis(redisClient *redis.Client) {
	RedisClient = redisClient
}
