package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/jessevdk/go-flags"
	rutil "github.com/jgilman1337/rod_util/pkg"
	"github.com/marusama/semaphore/v2"

	pkg "github.com/jgilman1337/chatbot_dl/internal"
	"github.com/jgilman1337/chatbot_dl/internal/cli"
	"github.com/jgilman1337/chatbot_dl/pkg/util"
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

	//Spawn a new web browser
	bopts := util.If(opts.Debug, rutil.DefaultBrowserOptsDbg(), rutil.DefaultBrowserOpts())
	bopts.Leakless = true
	browser, launcher, err := rutil.BuildSandboxless(bopts)
	if err != nil {
		log.Fatalf("Failed to launch Rod browser; reason: %s", err)
	}
	defer rutil.RodFree(browser, launcher)

	/*
		Create a page pool and semaphore, each with the size equal to that of the threads option
		The semaphore guarantees that only n number of threads are archived at once
		It also allows for the queuing up of overflow threads to archive once those before it finish
	*/
	pool := rod.NewPagePool(opts.Threads)
	defer pool.Cleanup(func(p *rod.Page) { p.MustClose() })
	sem := semaphore.New(opts.Threads)

	//TODO: temp stuff begin

	// Create a page if needed
	/*
		create := func() (*rod.Page, error) {
			// We use MustIncognito to isolate pages with each other
			return browser.MustPage(), nil
		}
	*/

	/*
		yourJob := func(id int, pp *rod.Pool[rod.Page]) {
			fmt.Printf("[%d] Spawned new thread\n", id)

			page, err := pp.Get(create)
			if err != nil {
				log.Fatalf("error while acquiring page: %s", err)
			}
			defer pp.Put(page)

			page.MustNavigate("http://example.com")
			page.MustWaitDOMStable()
			fmt.Printf("[%d] %s\n", id, page.MustInfo().Title)
		}
	*/

	//TODO: temp stuff end

	//Setup a waitgroup and semaphore
	wg := sync.WaitGroup{}

	//Archive each thread given in the positional arguments
	for i, url := range opts.Positional.URLs {
		//Skip URLs that do not correspond to a valis service
		_, tid, err := pkg.PickService(url)
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
			defer sem.Release(1) //Very important to prevent starvation

			//Run the worker
			//yourJob(i+1, &pool)
			if err := cli.RunPoolWorker(i+1, tid, &pool, *opts); err != nil {
				log.Printf("An error occurred while archiving URL at position %d; %s\n", i+1, err)
			}
		}()

		//Wait a bit to prevent Rod from being overloaded; infinitely loads the page otherwise
		time.Sleep(500 * time.Millisecond)
	}

	//Wait for the goroutines to all complete before continuing
	//This blocks the main thread
	wg.Wait()

	fmt.Printf("%+v\n", opts)
}
