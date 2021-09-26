package service

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

type RedisService interface {
	Set(key string, data interface{}, time int) error
	Exists(key string) bool
	Get(key string) ([]byte, error)
	PutExpireTime(key string, time int)
	Delete(key string) (bool, error)
	LikeDeletes(key string) error
}

type redisService struct {
	rp *redis.Pool
}

// Set a key/value
func (r *redisService) Set(key string, data interface{}, time int) error {
	conn := r.rp.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return nil
}

//设置过期时间
func (r *redisService) PutExpireTime(key string, time int) {
	conn := r.rp.Get()
	defer conn.Close()
	conn.Do("EXPIRE", key, time)
}

// Exists check a key
func (r *redisService) Exists(key string) bool {
	conn := r.rp.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// Get get a key
func (r *redisService) Get(key string) ([]byte, error) {
	conn := r.rp.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func NewRedisService(rp *redis.Pool) RedisService {
	return &redisService{rp: rp}
}

// Delete delete a kye
func (r *redisService) Delete(key string) (bool, error) {
	conn := r.rp.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func (r *redisService) LikeDeletes(key string) error {
	conn := r.rp.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = r.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
