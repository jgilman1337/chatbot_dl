package perplexity

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"

	rutil "github.com/jgilman1337/rod_util/pkg"

	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
	"github.com/jgilman1337/chatbot_dl/pkg/util"
)

// Downloads a thread to a specified directory as either a Markdown or PDF file.
func dumpThread(b *rod.Browser, p *rod.Page, ctx context.Context, ttype c.ThreadType) ([]byte, error) {
	logger := c.LoggerFromCtx(ctx)

	lprefix := fmt.Sprintf("[dumpThread %s] ", ttype.NameFor())

	//Target the download button for the thread
	target := getSelector(ttype)
	exportBtn, err := rutil.SafeSelect(p, target)
	if err != nil || exportBtn == nil {
		return nil, fmt.Errorf("%w '%s'", ErrSelectorFailed, target)
	}
	logger.Debug(lprefix + "Found download button; clicking...")

	//Download the file
	var data []byte
	dlerr := rod.Try(func() {
		waitForDownload := waitDownload(b.Timeout(10 * time.Second))
		exportBtn.MustClick()
		logger.Debug(lprefix + "Clicked download button; waiting for completion...")
		data, err = waitForDownload()
		if err != nil {
			//Just log the error; Line 48 will handle the rest
			logger.Error(lprefix + err.Error())
		}
	})
	if dlerr != nil {
		return nil, dlerr
	}

	//Ensure something was actually downloaded
	if len(data) == 0 {
		return nil, ErrNoDownloadBytes
	}
	util.LogFmt(logger.Debug, "%s Done; downloaded %d bytes", lprefix, len(data))

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

// Download waiter function that doesn't panic on errors.
// Also fixes a nil pointer dereference. Known issue here: https://github.com/go-rod/rod/issues/916
func waitDownload(b *rod.Browser) func() ([]byte, error) {
	tmpDir := filepath.Join(os.TempDir(), "rod", "downloads")
	wait := b.WaitDownload(tmpDir)

	return func() ([]byte, error) {
		info := wait()
		if info != nil {
			path := filepath.Join(tmpDir, info.GUID)
			defer func() { _ = os.Remove(path) }()
			return os.ReadFile(path)
		} else {
			return nil, ErrNilDownloadWaiter
		}
	}
}
