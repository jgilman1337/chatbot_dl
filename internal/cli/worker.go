package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	rutil "github.com/jgilman1337/rod_util/pkg"

	"github.com/jgilman1337/chatbot_dl/pkg/service"
	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
	"github.com/jgilman1337/chatbot_dl/pkg/util"
)

// Holds the parameters for the worker.
type WorkerParams struct {
	//The ID of the task.
	ID int

	//The ID of the thread to archive.
	TID string

	//The service to use when archiving the thread.
	Srv service.ServiceWD

	//The page pool that this worker should use.
	//Pool *rod.Pool[rod.Page]

	//The parent browser instance.
	//Browser *rod.Browser
}

//TODO: better logger with formatting; takes in logger func and format args `func logFmt(lfunc func(msg string, args ...any), msg string, args ...any)`

// Runs an archival operation for a singular thread. This function also handles tasks like creation of the output directory.
func RunWorker(wp WorkerParams, opts Options) error {
	//Initialize slog and attach it to a new context
	lvl := util.If(opts.Verbose, slog.LevelDebug, slog.LevelInfo)
	logger, buf := util.NewBufSlogSH(&lvl)
	/*
		logger = logger.With(
			util.WorkerIDKey, wp.ID,
			util.ThreadIDKey, wp.TID,
			util.ServiceIdentKey, wp.Srv.Ident(),
		)
	*/
	ctx := context.Background()
	ctx = c.WithLogger(ctx, logger)
	logger.Info(
		fmt.Sprintf("Using service '%s' with thread ID '%s'", wp.Srv.Ident(), wp.TID),
	)

	//Create a writer buffer for slog for Rod
	wh := util.NewSlogWriterHook(logger, slog.LevelDebug)

	/*
		Create a new browser for this instance, at the cost of increased RAM usage, but ensures things run smoother

		Severe issues arouse when trying to use Rod's page pool. If someone has a solution, feel free to PR it. Old testing code was archived, and might be of use to the PRing party.
	*/
	bopts := util.If(opts.Debug, rutil.DefaultBrowserOptsDbg(), rutil.DefaultBrowserOpts())
	bopts.Leakless = true
	bopts.Logger = wh
	browser, launcher, err := rutil.BuildSandboxless(bopts)
	if err != nil {
		log.Fatalf("Failed to launch Rod browser; reason: %s", err)
	}
	defer rutil.RodFree(browser, launcher)

	//Create a stealth page to bypass CF Turnstile
	pg, err := c.CreateStealthPage(browser, nil, ctx)
	if err != nil {
		log.Fatalf("error while acquiring page: %s", err)
	}
	defer pg.Close()

	//Create a new stealth page and assign a timeout context window
	page := pg.Timeout(time.Duration(opts.Timeout) * time.Second)
	defer page.CancelTimeout()

	//Set options for the scraper
	//TODO maybe expand the options into an interface so binding is easier
	/*
		if wp.Srv.Ident() == pplx.Ident {
			pplxOpts := pplx.DefaultDLOpts()
			pplxOpts.Timeout = opts.Timeout
			//pplxOpts.Formats = opts.Formats
			ctx = pplx.WithOptions(ctx, &pplxOpts)
		}
	*/

	//Scrape the thread
	result, err := wp.Srv.Scrape(browser, page, ctx, wp.TID)
	if err != nil {
		if result != nil {
			//t.Log("Result is non-nil")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			//t.Log("Caught a context deadline exceeded error")
		}

		return err
	}

	//Create the necessary working directories
	wdir, err := createWorkdirs(wp.ID, opts.Workdir, opts.Flat)
	if err != nil {
		return err
	}

	//Add the logger buffer to the list of results
	result = append(result, c.Thread{
		Type:     c.Markdown,
		Filename: "log",
		Content:  buf.Bytes(),
	})

	//Save the results to the working directory
	for _, thread := range result {
		err := os.WriteFile(path.Join(wdir, thread.GetFilename()), thread.Content, 0644)
		if err != nil {
			return err
		}
	}

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
