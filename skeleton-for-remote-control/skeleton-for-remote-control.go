package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func main() {

	fmt.Println("Start")

	// A context for graceful shutdown (It is based on the signal package)
	//
	// NOTE
	// Use os.Interrupt to gracefully shutdown on Ctrl+C which is SIGINT
	// Use syscall.SIGTERM which is the usual signal for termination and
	// the default one (it can be modified) for docker containers, which is also used by kubernetes.
	gracefulShutdownContext, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Package errgroup provides synchronization, error propagation, and Context cancelation
	// for groups of goroutines working on subtasks of a common task.
	group, groupContext := errgroup.WithContext(gracefulShutdownContext)

	// Process section
	processContext, processCancel := context.WithCancel(context.TODO())
	group.Go(func() error {
		fmt.Println("Start - process section")

		for {

			select {
			case <-groupContext.Done():
				fmt.Printf("[Process section] Cancel the groupContext(%v)\n", groupContext.Err())
				return groupContext.Err()

			case <-processContext.Done():
				fmt.Printf("[Process section] Cancel the processContext(%v)\n", processContext.Err())
				processContext, processCancel = context.WithCancel(context.TODO())
				time.Sleep(100 * time.Millisecond)
				// return processContext.Err()
			}
		}
	})

	// Necessary?
	_, contorlCancel := context.WithCancel(context.TODO())
	defer contorlCancel()

	// Control section
	group.Go(func() error {
		fmt.Println("Start - control section")

		for {
			select {
			case <-groupContext.Done():
				fmt.Printf("[Control section] Cancel the groupContext(%v)\n", groupContext.Err())
				return groupContext.Err()

			case <-time.After(10 * time.Second):
				fmt.Println("[Control section] Suspending the process")
				processCancel()

				// case <-time.After(6 * time.Second):
				// 	fmt.Println("(TBD)Resuming the process")
			}
			fmt.Println("Looping")
		}
	})

	fmt.Printf("Wait until go routines are finished.\n")
	err := group.Wait()
	if err != nil {
		fmt.Printf("Error group: %v\n", err)
	}

	fmt.Println("Main done")
}
