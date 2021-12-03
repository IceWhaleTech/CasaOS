package random

import (
	"fmt"
	"testing"
)

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(6, true))
}
