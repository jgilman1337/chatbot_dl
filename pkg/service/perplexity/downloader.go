package perplexity

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-rod/rod"

	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
)

// Downloads a thread to a specified directory as either a Markdown or PDF file.
func dumpThread(b *rod.Browser, p *rod.Page, ctx context.Context, ttype c.ThreadType) ([]byte, error) {
	logger := c.LoggerFromCtx(ctx)

	lprefix := fmt.Sprintf("[dumpThread %s]", ttype.NameFor())

	//Target the download button for the thread
	target := getSelector(ttype)
	exportBtn, err := p.Element(target)
	if err != nil || exportBtn == nil {
		return nil, fmt.Errorf("failed to find download button; selector: '%s'", target)
	}
	logger.Debug(lprefix + " Found download button; clicking...")

	//Download the file
	var data []byte
	dlerr := rod.Try(func() {
		waitForDownload := b.MustWaitDownload()
		exportBtn.MustClick()
		logger.Debug(lprefix + " Clicked download button; waiting for completion...")
		data = waitForDownload()
	})
	if dlerr != nil {
		return nil, dlerr
	}

	//Ensure something was actually downloaded
	if len(data) == 0 {
		return nil, errors.New("nothing was downloaded or the output byte array is nil")
	}
	logger.Debug(fmt.Sprintf("%s Done; downloaded %d bytes", lprefix, len(data)))

	return data, nil
}

// Emits the correct CSS selector for the thread type to export.
func getSelector(t c.ThreadType) string {
	typ := ""
	switch t {
	case c.DOCX:
		typ = "[data-testid=\"thread-export-docx\"]"
	case c.Markdown:
		typ = "[data-testid=\"thread-export-md\"]"
	case c.PDF:
		typ = "[data-testid=\"thread-export-pdf\"]"
	}
	return typ
}
