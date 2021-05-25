package ds

import "testing"

func TestWaitGroup_Add(t *testing.T) {
	var wg WaitGroup
	wg.Add(1)
}
