package cmd

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr:     config.RedisOption.Addr,     //"localhost:6379"
	Password: config.RedisOption.Password, // no password set
	DB:       config.RedisOption.DB,       // use default DB
})

// read from redis cache
func (rs *RegexFactory) FromCache() bool {
	// if redis disabled, return
	if !config.RedisOption.Enable {
		return false
	}

	val, err := client.Get(ctx, rs.InputKey).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(err)
	}
	//var nrs RegexFactory
	err = json.Unmarshal([]byte(val), rs)
	if err != nil {
		panic(err)
	}
	return true
}

// save to redis as cache
func (rs *RegexFactory) ToCache() {
	// if redis disabled, return
	if !config.RedisOption.Enable {
		return
	}
	//byteArray, err := json.Marshal(rs)
	byteArray, err := json.MarshalIndent(rs, config.Prefix, config.Indent)
	if err != nil {
		panic(err)
	}
	err = client.Set(ctx, rs.InputKey, byteArray, 0).Err()
	if err != nil {
		panic(err)
	}
}
