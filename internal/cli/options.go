package cli

// Holds the options for the commandline interface.
type Options struct {
	TID     string //The ID of the thread to archive.
	WorkDir string //The directory to output the files to.
	Timeout uint   //The max time (in seconds) that the tool is allowed to take when scraping a thread.
}

// Returns the default options for the `Options` struct.
func DefaultOptions() Options {
	return Options{
		TID:     "",
		WorkDir: ".",
		Timeout: 20,
	}
}
