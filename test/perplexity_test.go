package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/jgilman1337/chatbot_dl/pkg/service/perplexity"
	rutil "github.com/jgilman1337/rod_util/pkg"
)

//TODO: detect fullscreen modals that sometimes show up; will fail if one pops up

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

	//Add options; slog will be auto-attached as the default instance
	ctx := context.Background()
	opts := perplexity.DefaultDLOpts()
	opts.Timeout = 20
	ctx = perplexity.WithOptions(ctx, &opts)

	//Scrape the thread
	result, err := perplexity.Scrape(browser, ctx, id)
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
