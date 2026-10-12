package redis

import (
	"backend-rms/core/config"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	BASE_REDIS_PREFIX = "tb:cache:"
	RMS_REDIS_PREFIX  = "rms:cache:"
)

type RedisClient struct {
	client  *redis.Client
	context context.Context
	prefix  string
}

var BaseRedisCli *RedisClient
var RMSRedisCli *RedisClient
var QuoteRedisCli *RedisClient
var PositionRedisCli *RedisClient

func InitBaseRedis(conf config.RedisConf) {
	BaseRedisCli = initRedis(conf, BASE_REDIS_PREFIX)
}

func InitRMSRedis(conf config.RedisConf) {
	RMSRedisCli = initRedis(conf, RMS_REDIS_PREFIX)
}

func InitQuoteRedis(conf config.RedisConf) {
	QuoteRedisCli = initRedis(conf, "")
}

func InitPositionRedis(conf config.RedisConf) {
	PositionRedisCli = initRedis(conf, "")
}

func initRedis(conf config.RedisConf, prefix string) *RedisClient {
	if conf == (config.RedisConf{}) {
		panic("init redis failed, please confirm the Redis config is set correctly")
	}
	context := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return &RedisClient{client: client, context: context, prefix: prefix}
}

func (rc RedisClient) Exists(keys ...string) (int64, error) {
	for i, key := range keys {
		keys[i] = rc.prefix + key
	}

	result, err := rc.client.Exists(rc.context, keys...).Result()
	return result, errors.Wrap(err, "redis: get the key is exists failed")
}

func (rc RedisClient) Expire(key string, expiration time.Duration) (bool, error) {
	key = rc.prefix + key
	result, err := rc.client.Expire(rc.context, key, expiration).Result()

	return result, errors.Wrap(err, "redis: get the key's expire failed")
}

func (rc RedisClient) TTL(key string) (time.Duration, error) {
	key = rc.prefix + key
	result, err := rc.client.TTL(rc.context, key).Result()

	return result, errors.Wrap(err, "redis: get the key's ttl failed")
}

func (rc RedisClient) Get(key string) (string, error) {
	key = rc.prefix + key
	result, err := rc.client.Get(rc.context, key).Result()

	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	return result, errors.Wrap(err, "redis: get the key's value failed")
}

func (rc RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	key = rc.prefix + key
	err := rc.client.Set(rc.context, key, value, expiration).Err()

	return errors.Wrap(err, "redis: set the key failed")
}

func (rc RedisClient) SetEX(key string, value interface{}, expiration time.Duration) error {
	key = rc.prefix + key
	_, err := rc.client.SetEX(rc.context, key, value, expiration).Result()

	return errors.Wrap(err, "redis: setEX the key failed")
}

func (rc RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	key = rc.prefix + key
	result, err := rc.client.SetNX(rc.context, key, value, expiration).Result()

	return result, errors.Wrap(err, "redis: setNX the key failed")
}

func (rc RedisClient) HKeys(key string) ([]string, error) {
	key = rc.prefix + key
	result, err := rc.client.HKeys(rc.context, key).Result()

	return result, errors.Wrap(err, "redis: get the hash Hkeys failed")
}

func (rc RedisClient) HLen(key string) (int64, error) {
	key = rc.prefix + key
	result, err := rc.client.HLen(rc.context, key).Result()

	return result, errors.Wrap(err, "redis: get the hash hlen failed")
}

func (rc RedisClient) HGet(key, field string) (string, error) {
	key = rc.prefix + key
	result, err := rc.client.HGet(rc.context, key, field).Result()

	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	return result, errors.Wrap(err, "redis: get hash key's field value failed")
}

// HMGet returns the values for the specified fields in the hash stored at key.
// It returns an interface{} to distinguish between empty string and nil value.
func (rc RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	key = rc.prefix + key
	result, err := rc.client.HMGet(rc.context, key, fields...).Result()

	return result, errors.Wrap(err, "redis: get hash key's fields value failed")
}

func (rc RedisClient) HGetAll(key string) (map[string]string, error) {
	key = rc.prefix + key
	result, err := rc.client.HGetAll(rc.context, key).Result()

	return result, errors.Wrap(err, "redis: get hash key's value failed")
}

// HSet accepts values in following formats:
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (rc RedisClient) HSet(key string, values ...interface{}) (int64, error) {
	key = rc.prefix + key
	result, err := rc.client.HSet(rc.context, key, values...).Result()

	return result, errors.Wrap(err, "redis: set hash key value failed")
}

func (rc RedisClient) Del(keys ...string) (int64, error) {
	for i, key := range keys {
		keys[i] = rc.prefix + key
	}

	result, err := rc.client.Del(rc.context, keys...).Result()
	return result, errors.Wrap(err, "redis: delete key failed")
}
