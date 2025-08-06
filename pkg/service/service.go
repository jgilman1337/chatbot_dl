package service

import (
	"context"
	"strings"

	"github.com/go-rod/rod"

	c "github.com/jgilman1337/chatbot_dl/pkg/service/common"
)

// Represents a service that uses a Rod-powered web driver to scrape a website.
type ServiceWD interface {
	//Builds a link from a thread ID.
	BuildLink(tid string) string

	//Returns the identity (name) for the service. Ideally, this should be the domain name.
	Ident() string

	//Determines if a link is valid for the given service.
	IsValidLink(l string) bool

	//The scraper function to call.
	Scrape(b *rod.Browser, p *rod.Page, ctx context.Context, tid string) (sres []c.Thread, serr error)

	//Gets the base URL for the selected service.
	Stem() string
}

// Default implementation of `ServiceWD.BuildLink`.
func BuildLink(l string, stem string) string {
	if !strings.HasSuffix(stem, "/") {
		stem += "/"
	}
	return stem + l
}

// Default implementation of `ServiceWD.IsValidLink`.
func IsValidLink(l string, stem string) bool {
	return strings.HasPrefix(l, stem)
}
