package gredis

import (
	"github.com/gomodule/redigo/redis"
	"oasis/model"
	"time"
)

func GetRedisConn(m *model.RedisModel) *redis.Pool {
	redisConn := &redis.Pool{
		MaxIdle:     m.MaxIdle,
		MaxActive:   m.MaxActive,
		IdleTimeout: m.IdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", m.Host)
			if err != nil {
				return nil, err
			}
			if m.Password != "" {
				if _, err := c.Do("AUTH", m.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return redisConn
}
