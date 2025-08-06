package pkg

import (
	"fmt"
	"maps"
	"net/url"

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
func PickService(u string) (*service.ServiceWD, error) {
	//Parse the incoming URL
	urll, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	//Pick an appropriate service
	service, ok := Registry[urll.Host]
	if !ok {
		return nil, fmt.Errorf("unsupported service. The list of supported services are: %v", maps.Keys(Registry))
	}

	return &service, nil
}
