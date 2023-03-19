package file

import (
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

func TestNameAccumulation(t *testing.T) {
	goleak.VerifyNone(t)

	fmt.Println("aaa")
	a := NameAccumulation("/mnt/test_1_1", "/")
	fmt.Println(a)
}
