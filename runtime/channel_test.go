package runtime

import (
	"fmt"
	"testing"
)

func TestGoChan_getDataFromBuf(t *testing.T) {
	s1 := []int{1, 2, 3, 0, 0}
	target := make([]int, 5)
	copy(target, s1[1:])
	fmt.Println(target)
}
