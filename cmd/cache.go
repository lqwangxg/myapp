package cmd

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr:     config.RedisOption.Addr,     //"localhost:6379"
	Password: config.RedisOption.Password, // no password set
	DB:       config.RedisOption.DB,       // use default DB
})

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

func (rs *RegexFactory) ToCache() {
	// if redis disabled, return
	if !config.RedisOption.Enable {
		return
	}
	byteArray, err := json.Marshal(rs)
	//byteArray, err := json.MarshalIndent(rs, config.Prefix, config.Indent)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("MarshalIndent: %s", string(byteArray))

	// var nrs RegexFactory
	// err = json.Unmarshal(byteArray, &nrs)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// fmt.Printf("Unmarshal: %+v", nrs)
	err = client.Set(ctx, rs.InputKey, byteArray, 0).Err()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("setValue: key=%s,  rule=%s \n", key, value)
}

func JsonStringAutoDecode() func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
	return func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
		if rf != reflect.String || rt == reflect.String {
			return data, nil
		}
		raw := data.(string)
		if raw != "" && json.Valid([]byte(raw)) {
			var m interface{}
			err := json.Unmarshal([]byte(raw), &m)
			return m, err
		}
		return data, nil
	}
}
