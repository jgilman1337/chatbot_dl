package pkg

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/jgilman1337/chatbot_dl/pkg/service"
	"github.com/jgilman1337/chatbot_dl/pkg/service/perplexity"
)

// Holds the default services list in the form of a map.
var Registry = make(map[string]service.ServiceWD, 0)

// This function is automatically called when the package is loaded.
func init() {
	//Services list
	services := []service.ServiceWD{
		perplexity.PplxScraper{},
	}

	//Add the services to the registry
	for _, service := range services {
		Registry[service.Ident()] = service
	}
}

// Picks an appropriate service to handle the URL from the registry based on the URL stem.
func PickService(u string) (service.ServiceWD, string, error) {
	//Parse the incoming URL
	urll, err := url.Parse(u)
	if err != nil {
		return nil, "", err
	}

	//Pick an appropriate service
	tld := getRegisteredDomain(urll.Host)
	service, ok := Registry[tld]
	if !ok {
		return nil, "", fmt.Errorf(
			"unsupported service '%s'. The list of supported services are: %v",
			tld,
			strings.Join(slices.Collect(maps.Keys(Registry)), ", "),
		)
	}

	//Ensure the thread URL prefix matches what is expected from the service
	if !strings.HasPrefix(u, service.Stem()) {
		return nil, "", fmt.Errorf(
			"unrecognized stem for service '%s'. The URL must have the following stem: %s",
			service.Ident(),
			service.Stem(),
		)
	}

	return service, service.GetThreadID(urll), nil
}

// Gets the top-level domain name and TLD from a hostname. Only works for basic TLDs and not compound ones like `.co.uk`, `.gov.uk`, etc.
func getRegisteredDomain(host string) string {
	parts := strings.Split(host, ".")
	l := len(parts)
	if l >= 2 {
		return parts[l-2] + "." + parts[l-1]
	}
	return host
}
