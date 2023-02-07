package file

import (
	"fmt"
	"testing"
)

func TestNameAccumulation(t *testing.T) {
	fmt.Println("aaa")
	a := NameAccumulation("/mnt/test_1_1", "/")
	fmt.Println(a)
}
