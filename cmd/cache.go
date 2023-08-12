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
func (rs *Regex) FromCache(input string) bool {
	// if redis disabled, return
	if !rs.Cache {
		return false
	}

	rs.CacheKey = rs.Hashsum(input)
	val, err := client.Get(ctx, rs.CacheKey).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(err)
	}
	//var nrs RegexFactory
	err = json.Unmarshal([]byte(val), &rs.Result)
	if err != nil {
		panic(err)
	}

	return true
}

// save to redis as cache
func (rs *Regex) ToCache() {
	// if redis disabled, return
	if !rs.Cache {
		return
	}
	//byteArray, err := json.Marshal(rs)
	byteArray, err := json.MarshalIndent(rs.Result, config.Prefix, config.Indent)
	if err != nil {
		panic(err)
	}
	err = client.Set(ctx, rs.CacheKey, byteArray, 0).Err()
	if err != nil {
		panic(err)
	}
}
