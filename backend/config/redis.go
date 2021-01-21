package config

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

func CreateRedisClient(conf Configuration) {
	opt, err := redis.ParseURL(conf.REDIS_HOST)
	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(opt)
	Redis = redis
}
