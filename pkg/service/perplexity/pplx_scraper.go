package perplexity

import (
	"net/url"
	"strings"

	"github.com/jgilman1337/chatbot_dl/pkg/service"
)

const (
	//The identity string of this scraper.
	Ident = "perplexity.ai"
)

// Enforces compliance with the IConfig interface.
var _ service.ServiceWD = (*PplxScraper)(nil)

// Represents a scraper for Perplexity.ai.
type PplxScraper struct {
}

// Implements the BuildLink() function from ServiceWD.
func (s PplxScraper) BuildLink(tid string) string {
	return service.BuildLink(tid, s.Stem())
}

// Implements the GetThreadID() function from ServiceWD.
func (s PplxScraper) GetThreadID(u *url.URL) string {
	pieces := strings.Split(u.Path, "/")
	return pieces[len(pieces)-1] //Last item in the pieces array
}

// Implements the Ident() function from ServiceWD.
func (s PplxScraper) Ident() string {
	return Ident
}

// Implements the IsValidLink() function from ServiceWD.
func (s PplxScraper) IsValidLink(l string) bool {
	return service.IsValidLink(l, s.Stem())
}

// Implements the Stem() function from ServiceWD.
func (s PplxScraper) Stem() string {
	return "https://www.perplexity.ai/search/"
}
