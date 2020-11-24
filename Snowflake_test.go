package Snowflake

import (
	"testing"
)

func Test_Generate(t *testing.T) {
	worker, err := NewWorker(1)
	checkMap := make(map[int64]struct{})
	if err != nil {
		t.Error(err)
	}

	ch := make(chan int64, 10)

	for count := 1; count <= 10000; count++ {
		go func(ch chan int64) {
			ch <- worker.Generate()
		}(ch)
	}

	for count := 1; count <= 10000; count++ {
		id := <-ch
		if _, isExist := checkMap[id]; isExist {
			t.Error("already exist id")
		} else {
			checkMap[id] = struct{}{}
		}
	}
}
