package threadpool

import (
	"fmt"
	"time"
)

type ThreadPool struct {
	workerCount int
	taskQueue   chan interface{}
}

func playWithChannels() {
	ch := make(chan int)

	go func() {
		counter := 0
		for counter < 10 {
			time.Sleep(time.Second)
			counter++
			ch <- counter
		}
		close(ch)
	}()

	for x := range ch {
		fmt.Println(x)
	}
}
