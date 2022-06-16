package main

import (
	"fmt"
	"time"
)

func main() {
	duration, _ := time.ParseDuration("10s")
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	fmt.Println(time.Now())
	fmt.Println("Sleep 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("Wake up")

	for {
		select {
		case t := <-ticker.C:
			fmt.Printf("Tick at %v\n", t)

			fmt.Println("Sleep 3 seconds")
			time.Sleep(3 * time.Second)
			fmt.Println("Wake up")
		}
	}
}
