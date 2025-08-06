package perplexity

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	rutil "github.com/jgilman1337/rod_util/pkg"

	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
	"github.com/jgilman1337/chatbot_dl/pkg/util"
)

// Implements the Scrape() function from ServiceWD.
func (s PplxScraper) Scrape(b *rod.Browser, p *rod.Page, ctx context.Context, tid string) (sres []c.Thread, serr error) {
	logger := c.LoggerFromCtx(ctx)
	opts := OptsFromCtx(ctx)
	result := make([]c.Thread, 0)

	//Create a page if it doesn't yet exist
	createdPage := p == nil
	if createdPage {
		logger.Debug("Attempting to self-spawn a new page...")

		np, err := c.CreateStealthPage(b, opts.Device, opts.Timeout, ctx)
		if err != nil {
			return nil, err
		}
		p = np

		logger.Debug("Spawned a new stealth page successfully")
	}

	//Defer the page close operation for later, but only if the function created the page
	caughtTimeoutErr := false
	defer func() {
		if !createdPage {
			return
		}

		logger.Debug("Attempting to kill the self-spawned page...")

		err := p.Close()
		if res := c.HandleErr(logger,
			c.NewErrHandlerParams(
				"Error while closing page",
				"page close",
				err,
				&caughtTimeoutErr,
				&serr,
			),
		); res != c.EH_OK {
			return
		}

		logger.Debug("Successfully killed the self-spawned page")
	}()

	//TODO: setup network monitoring stuff from old tests

	//Navigate to the target thread
	if err := p.Navigate(s.BuildLink(tid)); err != nil {
		if res := c.HandleErr(logger,
			c.NewErrHandlerParams(
				"Error while navigating to the target",
				"navigate to target",
				err,
				&caughtTimeoutErr,
				&serr,
			),
		); res != c.EH_OK {
			return
		}
	}
	logger.Info("Successfully loaded the webpage; waiting for the DOM to stabilize")

	//Wait for the page to load before continuing
	if err := p.WaitDOMStable(time.Second, 0); err != nil {
		if res := c.HandleErr(logger,
			c.NewErrHandlerParams(
				"Error while waiting on the DOM",
				"wait on DOM",
				err,
				&caughtTimeoutErr,
				&serr,
			),
		); res != c.EH_OK {
			return
		}
	}
	logger.Info("Successfully waited on the DOM; beginning archival process")

	//Download the threads as Markdown and PDF
	formats := opts.Formats
	for i, format := range formats {
		lprefix := fmt.Sprintf("[dl %s; %d/%d]", format.NameFor(), i+1, len(formats))

		//Close any modals that could interfere with the archival
		dismissModal(p, ctx)

		//Wait 1-3 seconds before opening the drawer
		logger.Info(lprefix + " Waiting to open thread download drawer...")
		time.Sleep(time.Millisecond * time.Duration(
			opts.DLWaitMin+
				rand.N(opts.DLWaitMax)),
		)

		//Target the drawer button for the thread
		logger.Info(lprefix + " Opening thread download drawer...")
		drawerSelector := "[data-testid=\"thread-dropdown-menu\"]"
		drawerBtn, err := rutil.SafeSelect(p, drawerSelector)
		if err != nil || drawerBtn == nil {
			if res := c.HandleErr(logger,
				c.NewErrHandlerParams(
					fmt.Sprintf("%s Failed to find drawer button (skipped); selector: '%s'", lprefix, drawerSelector),
					lprefix+" find drawer button",
					err,
					&caughtTimeoutErr,
					&serr,
				),
			); res != c.EH_OK {
				if opts.AbortOnArchiveFailure {
					return
				} else {
					continue
				}
			}
		}

		//Click the thread download drawer
		if err := drawerBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
			if res := c.HandleErr(logger,
				c.NewErrHandlerParams(
					fmt.Sprintf("%s Failed to click drawer button (skipped); selector: '%s'", lprefix, drawerSelector),
					lprefix+" click drawer button",
					err,
					&caughtTimeoutErr,
					&serr,
				),
			); res != c.EH_OK {
				if opts.AbortOnArchiveFailure {
					return
				} else {
					continue
				}
			}
		}
		logger.Info(lprefix + " Thread download drawer clicked. Preparing to archive thread...")
		if opts.DrawerInteractDelay > 0 {
			time.Sleep(time.Millisecond * time.Duration(opts.DrawerInteractDelay))
		}
		//Archive the thread
		logger.Info(lprefix + " Archiving thread...")
		bytes, err := dumpThread(b, p, ctx, format)
		if err != nil {
			if res := c.HandleErr(logger,
				c.NewErrHandlerParams(
					"Failed to archive thread",
					lprefix+" archive thread",
					err,
					&caughtTimeoutErr,
					&serr,
				),
			); res != c.EH_OK {
				if opts.AbortOnArchiveFailure {
					return
				} else {
					continue
				}
			}
		}

		//Add the bytes to the output
		archive := c.Thread{
			Type:     format,
			Filename: util.Timestamp() + "-" + getPageTitle(p, ctx),
			Content:  bytes,
		}
		result = append(result, archive)
		fmt.Printf("Filename is '%s'\n", archive.Filename)

		logger.Info(lprefix + " Successfully archived thread")
	}

	return result, nil
}

// Gets the title of the page, truncating the title if it's too long.
func getPageTitle(p *rod.Page, ctx context.Context) string {
	logger := c.LoggerFromCtx(ctx)
	opts := OptsFromCtx(ctx)
	rawTitle := "perplexity_dl"

	errTxt := "Failed to get title; falling back to generic title \nCause: "

	//Select the title element
	elem, err := rutil.SafeSelect(p, "h1.group\\/query")
	if err != nil {
		logger.Warn(errTxt + err.Error())
		return rawTitle //Errors are ignored
	}
	if elem == nil {
		logger.Warn(errTxt + "<Nil title element>")
		return rawTitle //Errors are ignored
	}

	//Get the title of the thread
	t, err := elem.Text()
	if err != nil {
		logger.Warn(errTxt + err.Error())
		return rawTitle //Errors are ignored
	}
	rawTitle = t

	//Truncate to x chars
	runes := []rune(rawTitle)
	if len(runes) > opts.TNLen {
		rawTitle = string(runes[0:opts.TNLen]) + "..."
	}

	return rawTitle
}

// Dismisses modals if any appear.
func dismissModal(p *rod.Page, ctx context.Context) bool {
	logger := c.LoggerFromCtx(ctx)

	//Find the modal by selector
	modalSelector := "[data-testid=\"close-modal\"]"
	closeBtn, err := rutil.SafeSelect(p, modalSelector)
	if err != nil || closeBtn == nil {
		return false
	}

	//Click the close button on the modal
	if err := closeBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return false
	}

	logger.Debug(
		fmt.Sprintf("Automatically found and dismissed a full-screen modal; selector: '%s'", modalSelector),
	)
	return true
}
