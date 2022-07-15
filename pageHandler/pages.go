package pageHandler

var providers map[string]PageProvider

func GetProviders(cacheTemplates bool, dataStorage string) map[string]PageProvider {
	if providers == nil {
		providers = make(map[string]PageProvider)
		//Add the providers in the pages sub package
	}
	return providers
}
