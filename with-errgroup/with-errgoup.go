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

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {

		for {
			// NOTE - Default Selection
			// The default case in a select is run if no other case is ready.
			// Use a default case to try a send or receive without blocking:

			select {
			case <-gCtx.Done():
				fmt.Println("Break the loop")
				return nil
			case <-time.After(1 * time.Second):
				fmt.Println("Hello in a loop")
				// default:
				// 	fmt.Print(".")
			}

			fmt.Println("Bottom of for loop")
			// time.Sleep(100 * time.Millisecond)
		}
	})

	g.Go(func() error {
		<-gCtx.Done()
		fmt.Println("Caio in Second routine")
		return nil
		// for {
		// 	select {
		// 	case <-gCtx.Done():
		// 		fmt.Println("Break the loop")
		// 		return nil
		// 	case <-time.After(1 * time.Second):
		// 		fmt.Println("Ciao in a loop")
		// 	}
		// }
	})

	err := g.Wait()
	if err != nil {
		fmt.Println("Error group: ", err)
	}
	fmt.Println("Main done")

}
