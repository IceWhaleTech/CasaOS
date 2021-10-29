package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func Init() *cache.Cache {
	return cache.New(5*time.Minute, 60*time.Second)
}
