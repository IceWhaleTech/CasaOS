package sign

import (
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS/pkg/sign"
)

var once sync.Once
var instance sign.Sign

func Sign(data string) string {

	return NotExpired(data)

}

func WithDuration(data string, d time.Duration) string {
	once.Do(Instance)
	return instance.Sign(data, time.Now().Add(d).Unix())
}

func NotExpired(data string) string {
	once.Do(Instance)
	return instance.Sign(data, 0)
}

func Verify(data string, sign string) error {
	once.Do(Instance)
	return instance.Verify(data, sign)
}

func Instance() {
	instance = sign.NewHMACSign([]byte("token"))
}
