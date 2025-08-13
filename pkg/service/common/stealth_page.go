package common

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/stealth"

	rutil "github.com/jgilman1337/rod_util/pkg"
)

// Creates a stealthy page to aid in bypassing bot bans.
func CreateStealthPage(b *rod.Browser, d *devices.Device, ctx context.Context) (*rod.Page, error) {
	logger := LoggerFromCtx(ctx)

	//Create the page
	var p *rod.Page
	err := rod.Try(func() {
		p = stealth.MustPage(b)

		//Spoof the user agent
		dev := d
		if dev == nil {
			rdev := rutil.PickRandMobileDevice()
			dev = &rdev
		}
		p.MustEmulate(*dev)

		logger.Info(fmt.Sprintf("Using fake device '%s'; user agent: '%s'", dev.Title, dev.UserAgent))
	})

	//Log any errors that occurred
	if err != nil {
		LogErr(err, NewErrHandlerParams(
			"Error during page creation; cannot continue",
			"page creation",
			logger,
		))

		return nil, err
	}

	return p, nil
}
