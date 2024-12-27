package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type ThreadPool struct {
	workerCount int
	taskQueue   chan interface{}
}

func singleChannelWithMultipleListeners() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		counter := 0
		for counter < 10 {
			time.Sleep(time.Second)
			counter++
			ch <- counter
		}
		wg.Done()
		close(ch)
	}()
	for i := 0; i < runtime.NumCPU(); i++ {

		wg.Add(1)
		go func() {
			for x := range ch {
				fmt.Println(x * i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func multipleChannelWithMultipleListeners() {
	listeners := make([]chan int, 10)
	wg := sync.WaitGroup{}

	go func() {
		counter := 0
		for counter < 10 {
			time.Sleep(time.Second)
			counter++
			for _, ch := range listeners {
				ch <- counter
			}
		}
		wg.Done()
		for _, ch := range listeners {

			close(ch)
		}
	}()
	for i := 0; i < runtime.NumCPU(); i++ {

		wg.Add(1)
		ch := make(chan int)
		listeners = append(listeners, ch)
		go func(ch chan int) {
			for x := range ch {
				fmt.Println(x * i)
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
}
