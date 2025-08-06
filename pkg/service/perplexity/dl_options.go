package perplexity

import (
	"github.com/go-rod/rod/lib/devices"
	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
)

// Holds the options for a Perplexity thread archive operation.
type DLOpts struct {
	Timeout uint           //Maximum time allowed for the download operation; 0 to disable.
	Formats []c.ThreadType //The list of formats to archive.

	AbortOnArchiveFailure bool //Whether to completely exit if a format failed to be downloaded.

	DLWaitMin           uint //The minimum threshold for the page event runners (ms).
	DLWaitMax           uint //The minimum threshold for the page event runners (ms).
	DrawerInteractDelay uint //The time (ms) to wait before interacting with the archival drawer.

	TNLen int //The maximum length of a downloaded thread name.

	Device *devices.Device //The device to emulate when scraping. Leave as `nil` for a random selection.
}

// Returns the default options for a download options struct.
func DefaultDLOpts() DLOpts {
	return DLOpts{
		Timeout: 30,
		Formats: []c.ThreadType{c.Markdown, c.PDF},

		AbortOnArchiveFailure: true,

		DLWaitMin:           500,
		DLWaitMax:           2500,
		DrawerInteractDelay: 250,

		TNLen: 30,

		Device: nil,
	}
}
