package pageHandler

import "golang.captainalm.com/cityuni-webserver/pageHandler/pages"

var providers map[string]PageProvider

func GetProviders(cacheTemplates bool, dataStorage string) map[string]PageProvider {
	if providers == nil {
		providers = make(map[string]PageProvider)
		providers["/test.go"] = pages.NewTestPage() //Test Page
		//Add the providers in the pages sub package
	}
	return providers
}
