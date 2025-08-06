package cli

import (
	"github.com/jessevdk/go-flags"
)

// Holds the options for the commandline interface.
type Options struct {
	Formats []string `short:"f" long:"formats" description:"The list of formats to archive" default:"md pdf"`
	Timeout uint     `short:"t" long:"timeout" description:"The max time (in seconds) that the tool is allowed to take when scraping a thread" default:"20"`

	Workdir string `short:"d" long:"workdir" description:"The directory to place archived threads and logs in" default:"."`
	Flat    bool   `short:"F" long:"flat" description:"Whether to put all archived threads in the same folder or in sub-folders, numbered by their index"`

	Verbose bool `short:"v" long:"verbose" description:"Show verbose logs"`

	Threads uint `short:"T" long:"threads" description:"The number of goroutines to use to archive chatbot threads" default:"5" validate:"gte=1,lte=10"`

	Positional Positional `positional-args:"true"`
}

// Holds positional arguments.
type Positional struct {
	URLs []string `description:"The list of URLs to process" required:"true" validate:"required,urlslice"`
}

// Returns the default options for the `Options` struct.
func DefaultOptions() *Options {
	//Hack for defaults generation since creasty/defaults has a different array syntax to go-flags
	opts := Options{}
	if _, err := flags.ParseArgs(&opts, []string{""}); err != nil {
		panic(err) //This shouldn't ever be hit
	}
	return &opts
}
