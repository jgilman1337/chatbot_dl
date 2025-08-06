package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
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

	fmt.Printf("%+v\n", opts)
}
