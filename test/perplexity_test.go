package test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/devices"
	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
	"github.com/jgilman1337/chatbot_dl/pkg/service/perplexity"
	"github.com/jgilman1337/chatbot_dl/pkg/util"
	rutil "github.com/jgilman1337/rod_util/pkg"
)

func TestPerplexityBasic(t *testing.T) {
	//Launch a new browser with default options, and connect to it
	bopts := rutil.DefaultBrowserOptsDbg()
	bopts.Leakless = true
	browser, launcher, err := rutil.BuildSandboxless(bopts)
	if err != nil {
		panic(fmt.Errorf("Failed to launch Rod browser; reason: %w", err))
	}
	defer rutil.RodFree(browser, launcher)

	id := "ai-amplifies-false-memories-9iZN5JuFT5.9asR1Ntf._A"

	//Create a basic slog handler
	handler := util.NewWorkerHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	//Add options
	ctx := context.Background()
	opts := perplexity.DefaultDLOpts()
	opts.DLWaitMax = 750
	ctx = c.WithLogger(ctx, logger)
	ctx = perplexity.WithOptions(ctx, &opts)

	//Create the page
	pg, err := c.CreateStealthPage(
		browser, &devices.Nexus7, ctx,
	)
	if err != nil {
		t.Fatal(err)
	}

	//Temporary page wrapper with timeout for this job only
	p := pg.Timeout(20 * time.Second)
	defer p.CancelTimeout()

	//Scrape the thread
	pplx := perplexity.PplxScraper{}
	result, err := pplx.Scrape(browser, p, ctx, id)
	if err != nil {
		if result != nil {
			t.Log("Result is non-nil")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			t.Log("Caught a context deadline exceeded error")
		}
		t.Fatal(err)
	}

	//Archive results to temp directory
	tmpDir, err := os.MkdirTemp("", "perplexity_dl-*")
	if err != nil {
		t.Fatal(err)
	}
	for _, thread := range result {
		err := os.WriteFile(path.Join(tmpDir, thread.GetFilename()), thread.Content, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("Dumped threads to '%s'\n", tmpDir)

	//Auto-open the folder if on Linux; silently fail if this doesn't work
	if runtime.GOOS == "linux" {
		exec.Command("xdg-open", tmpDir).Run()
	}
}
