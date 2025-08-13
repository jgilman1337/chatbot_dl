package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/marusama/semaphore/v2"

	pkg "github.com/jgilman1337/chatbot_dl/internal"
	"github.com/jgilman1337/chatbot_dl/internal/cli"
)

func main() {
	//Parse args and check for errors
	opts, err := cli.ParseArgs(os.Args[1:]) //argv[0] is the prog name; strip it
	if err != nil {
		//Echo the error if it's not coming from go-flags
		if _, ok := err.(*flags.Error); !ok {
			fmt.Println(err)
		}

		os.Exit(1)
	}
	if opts == nil {
		os.Exit(0) //Possibly due to help text being shown
	}

	/*
		Create a semaphore, each with the size equal to that of the threads option
		The semaphore guarantees that only n number of threads are archived at once
		It also allows for the queuing up of overflow threads to archive once those before it finish
	*/
	sem := semaphore.New(opts.Threads)

	//Archive each thread given in the positional arguments
	wg := sync.WaitGroup{}
	for i, url := range opts.Positional.URLs {
		//Skip URLs that do not correspond to a valid service
		service, tid, err := pkg.PickService(url)
		if err != nil {
			log.Printf("Skipped URL at position %d; reason: %s\n", i+1, err)
			continue
		}

		//Begin the archival process inside a goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()

			/*
				Limit how many concurrent goroutines can archive threads at a time
				If the semaphore is full, then the current goroutine will wait for the semaphore to yield a new slot

				Below this point and up until the end of this scope is a critical section guarded by the below semaphore
			*/
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Fatalf("Failed to acquire semaphore: %v", err)
			}
			defer func() {
				sem.Release(1) //Very important to prevent starvation
				fmt.Println("worker", i+1, "has released the semaphore")
			}()

			//Run the worker
			wp := cli.WorkerParams{
				ID:  i + 1,
				TID: tid,
				Srv: service,
			}
			if err := cli.RunWorker(wp, *opts); err != nil {
				log.Printf("An error occurred while archiving URL at position %d; %s\n", i+1, err)
			}

			fmt.Println("worker", i+1, "has finished")
		}()

		//Wait a bit to prevent Rod from being overloaded; infinitely loads the page otherwise
		time.Sleep(125 * time.Millisecond)
	}

	//Wait for the goroutines to all complete before continuing
	//This blocks the main thread
	wg.Wait()

	fmt.Printf("%+v\n", opts)
}
