package perplexity

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"

	rutil "github.com/jgilman1337/rod_util/pkg"

	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
	"github.com/jgilman1337/chatbot_dl/pkg/util"
)

var baseUrl = "https://www.perplexity.ai/search/"

// Scrapes a Perplexity thread using Rod.
func Scrape(b *rod.Browser, ctx context.Context, id string) (sres []c.Thread, serr error) {
	logger := c.LoggerFromCtx(ctx)
	opts := OptsFromCtx(ctx)
	result := make([]c.Thread, 0)

	//Add stealth to try to bypass Cloudflare Turnstile
	var p *rod.Page
	err := rod.Try(func() {
		p = stealth.MustPage(b)
		if opts.Timeout > 0 {
			p = p.Timeout(time.Duration(opts.Timeout) * time.Second)
		}

		//Spoof the user agent
		dev := opts.Device
		if dev == nil {
			rdev := rutil.PickRandMobileDevice()
			dev = &rdev
		}
		logger.Info(fmt.Sprintf("Using fake device '%s'; user agent: '%s'", dev.Title, dev.UserAgent))
		p.MustEmulate(*dev)
	})

	//TODO: setup network monitoring stuff from old tests

	//Defer the page close operation for later
	caughtTimeoutErr := false
	defer func() {
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
	}()

	//Ensure the page was created successfully
	if res := c.HandleErr(logger,
		c.NewErrHandlerParams(
			"Error during page creation; cannot continue",
			"page creation",
			err,
			&caughtTimeoutErr,
			&serr,
		),
	); res != c.EH_OK {
		return
	}

	//Navigate to the target thread
	if err := p.Navigate(baseUrl + id); err != nil {
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

		//Wait 1-3 seconds before opening the drawer
		logger.Info(lprefix + " Waiting to open thread download drawer...")
		time.Sleep(time.Millisecond * (time.Duration(opts.DLWaitMin) +
			time.Duration(rand.N(opts.DLWaitMax))),
		)

		//Target the drawer button for the thread
		logger.Info(lprefix + " Opening thread download drawer...")
		drawerSelector := "[data-testid=\"thread-dropdown-menu\"]"
		drawerBtn, err := p.Element(drawerSelector)
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
		time.Sleep(time.Millisecond * 250)

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

	//Get the title of the thread
	err := rod.Try(func() {
		rawTitle = p.MustElement("h1.group\\/query").MustText()
	})
	if err != nil {
		logger.Warn(fmt.Sprintf(
			"Failed to get title; falling back to generic title \nStacktrace: %s",
			err.Error(),
		))
		return rawTitle //Errors are ignored
	}

	//Truncate to x chars
	runes := []rune(rawTitle)
	if len(runes) > opts.TNLen {
		rawTitle = string(runes[0:opts.TNLen]) + "..."
	}

	return rawTitle
}
