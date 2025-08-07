package cli

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-rod/rod"
)

// Runs an archival operation for a singular thread. This function also handles tasks like creation of the output directory.
func RunPoolWorker(id int, tid string, pp *rod.Pool[rod.Page], opts Options) error {
	//Create the necessary working directories
	workdir, err := createWorkdirs(id, opts.Workdir, opts.Flat)
	if err != nil {
		return err
	}
	log.Printf("workdir: '%s'\n", workdir)

	return nil
}

// Creates the necessary working directories.
func createWorkdirs(id int, workdir string, flat bool) (string, error) {
	//Get the absolute version of the input path
	wdir, err := filepath.Abs(workdir)
	if err != nil {
		return "", err
	}

	//Create subdirs if the user opted to run this in non-flat mode
	if !flat {
		wdir = filepath.Join(wdir, strconv.Itoa(id)) //TODO: might want to use TID instead
	}

	//Create the base directory if it doesn't already exist
	if err := os.MkdirAll(wdir, 0755); err != nil {
		return "", err
	}

	return wdir, nil
}
