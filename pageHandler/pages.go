package pageHandler

var providers map[string]PageProvider

func GetProviders(cacheTemplates bool, dataStorage string, pageHandler *PageHandler) map[string]PageProvider {
	if providers == nil {
		providers = make(map[string]PageProvider)
		if pageHandler != nil {
			infoPage := newGoInfoPage(pageHandler, dataStorage, cacheTemplates)
			providers[infoPage.GetPath()] = infoPage //Go Information Page
		}

		//Add the providers in the pages sub package
	}
	return providers
}
