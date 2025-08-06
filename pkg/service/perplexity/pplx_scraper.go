package perplexity

import "github.com/jgilman1337/chatbot_dl/pkg/service"

// Enforces compliance with the IConfig interface.
var _ service.ServiceWD = (*PplxScraper)(nil)

// Represents a scraper for Perplexity.ai.
type PplxScraper struct {
}

// Implements the BuildLink() function from ServiceWD.
func (s PplxScraper) BuildLink(tid string) string {
	return service.BuildLink(tid, s.Stem())
}

// Implements the Ident() function from ServiceWD.
func (s PplxScraper) Ident() string {
	return "perplexity.ai"
}

// Implements the IsValidLink() function from ServiceWD.
func (s PplxScraper) IsValidLink(l string) bool {
	return service.IsValidLink(l, s.Stem())
}

// Implements the Stem() function from ServiceWD.
func (s PplxScraper) Stem() string {
	return "https://www.perplexity.ai/search/"
}
