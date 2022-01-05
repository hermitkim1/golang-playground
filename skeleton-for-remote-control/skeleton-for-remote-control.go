package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var MyCounter MyTestCounter = MyTestCounter{}

type MyTestCounter struct {
	i      int
	ctx    context.Context
	cancel func()
}

func (a *MyTestCounter) Start() {

	a.setup()

	if err := a.work(); err != nil {
		fmt.Printf("%#v\n", err)
	}

}

func (a *MyTestCounter) setup() {
	ctx, cancel := context.WithCancel(context.Background())

	a.ctx = ctx
	a.cancel = cancel

	a.i = 0
}

func (a *MyTestCounter) work() error {

	fmt.Println("Start - work")

	for {
		select {
		default:
			time.Sleep(250 * time.Millisecond)
			a.i++
			fmt.Println(a.i)
		case <-a.ctx.Done():
			fmt.Printf("closing work()\n")
			return a.ctx.Err()
		}
	}

}

func (a *MyTestCounter) tidyUp() {
	fmt.Println("Start - tidyUp")

	fmt.Println("End - tidyUp")
}

func (a *MyTestCounter) Stop() {
	a.cancel()
	a.tidyUp()
}

func process(command <-chan string, ctx context.Context, wg *sync.WaitGroup) {

	fmt.Println("Start - process section")

	defer wg.Done()

	go func() {
		for {
			cmd := <-command
			fmt.Println(cmd)
			switch cmd {
			case "suspend":
				MyCounter.Stop()
			case "resume":
				go MyCounter.Start()
			default:
				fmt.Println("Default ?")
			}
		}
	}()

	<-ctx.Done()

	// Something to do

	fmt.Println("End - process section")
}

func control(command chan string, ctx context.Context, wg *sync.WaitGroup) {

	fmt.Println("Start - control section")
	defer wg.Done()

	// An example of control

	time.Sleep(1 * time.Second)
	command <- "resume"

	time.Sleep(1 * time.Second)
	command <- "suspend"

	time.Sleep(3 * time.Second)
	command <- "resume"

	time.Sleep(1 * time.Second)
	command <- "suspend"

	<-ctx.Done()

	// Something to do

	fmt.Println("End - control section")
}

func main() {

	fmt.Println("Start")

	var wg sync.WaitGroup

	// A context for graceful shutdown (It is based on the signal package)
	//
	// NOTE
	// Use os.Interrupt to gracefully shutdown on Ctrl+C which is SIGINT
	// Use syscall.SIGTERM which is the usual signal for termination and
	// the default one (it can be modified) for docker containers, which is also used by kubernetes.
	gracefulShutdownContext, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-gracefulShutdownContext.Done()
		fmt.Println("Tasks before shutting down")

		// Add additional tasks here

		stop()
	}()

	// Package errgroup provides synchronization, error propagation, and Context cancelation
	// for groups of goroutines working on subtasks of a common task.
	// group, groupContext := errgroup.WithContext(gracefulShutdownContext)

	command := make(chan string)
	// Process section
	wg.Add(1)
	go process(command, gracefulShutdownContext, &wg)

	// Control section
	wg.Add(1)
	go control(command, gracefulShutdownContext, &wg)

	fmt.Printf("Wait until go routines are finished.\n")
	wg.Wait()
	// err := group.Wait()
	// if err != nil {
	// 	fmt.Printf("Error group: %v\n", err)
	// }

	fmt.Println("Main done")
}
