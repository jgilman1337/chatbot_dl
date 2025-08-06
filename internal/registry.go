package pkg

import (
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
