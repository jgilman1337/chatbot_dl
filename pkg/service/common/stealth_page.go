package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/stealth"

	rutil "github.com/jgilman1337/rod_util/pkg"
)

// Creates a stealthy page to aid in bypassing bot bans.
func CreateStealthPage(b *rod.Browser, d *devices.Device, timeout uint, ctx context.Context) (*rod.Page, error) {
	logger := LoggerFromCtx(ctx)

	//Create the page
	var p *rod.Page
	err := rod.Try(func() {
		p = stealth.MustPage(b) //TODO: add `.MustIncognito()`
		if timeout > 0 {
			p = p.Timeout(time.Duration(timeout) * time.Second)
		}

		//Spoof the user agent
		dev := d
		if dev == nil {
			rdev := rutil.PickRandMobileDevice()
			dev = &rdev
		}
		logger.Info(fmt.Sprintf("Using fake device '%s'; user agent: '%s'", dev.Title, dev.UserAgent))
		p.MustEmulate(*dev)
	})

	/*
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
	*/

	//Ensure the page was created successfully
	var oerr error
	firedCtxErr := false
	if res := HandleErr(logger,
		NewErrHandlerParams(
			"Error during page creation; cannot continue",
			"page creation",
			err,
			&firedCtxErr,
			&oerr,
		),
	); res != EH_OK {
		return nil, oerr
	}

	return p, nil
}
