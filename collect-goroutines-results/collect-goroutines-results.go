package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	contents := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	fmt.Printf("%+v\n", contents)
	fmt.Printf("%#v\n", contents)
	fmt.Printf("%v\n", contents)

	sampleChan := make(chan string)
	for _, line := range contents {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			sampleChan <- newSample(line)
		}(line)
	}

	go func() {
		fmt.Println("waiting")
		wg.Wait()
		close(sampleChan)
	}()

	for s := range sampleChan {
		fmt.Println(s)
	}

	fmt.Println("Finished")

}

func newSample(line string) string {
	return line
}
