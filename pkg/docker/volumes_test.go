package docker

import (
	"fmt"
	"testing"
)

func TestGetDir(t *testing.T) {
	fmt.Println(GetDir("", "config"))
}
