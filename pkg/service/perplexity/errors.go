package perplexity

import "errors"

var (
	ErrNilDownloadWaiter = errors.New("nil download waiter possibly due to cancelled context; cannot continue")
	ErrNoDownloadBytes   = errors.New("nothing was downloaded or the output byte array is nil")
	ErrSelectorFailed    = errors.New("failed to find element for selector")

	ErrNilBrowser = errors.New("nil browser instance; cannot continue")
	ErrNilPage    = errors.New("nil page instance; cannot continue")
)
