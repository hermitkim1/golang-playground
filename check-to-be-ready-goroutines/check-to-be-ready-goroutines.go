package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Start - main()")

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("Start - goroutine 1")
		for {
			<-time.After(1 * time.Second)
			fmt.Println("goroutine 1 - ticking every 1 sec")
		}
	}(&wg)
	time.Sleep(1 * time.Second)
	fmt.Println("Ready - goroutine 1")

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("Start - goroutine 2")
		for {
			<-time.After(2 * time.Second)
			fmt.Println("goroutine 2 - ticking every 2 sec")
		}
	}(&wg)
	time.Sleep(1 * time.Second)
	fmt.Println("Ready - goroutine 2")

	wg.Wait()

	fmt.Println("End - main()")
}
