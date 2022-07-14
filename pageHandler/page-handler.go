package pageHandler

import (
	"golang.captainalm.com/cityuni-webserver/conf"
	"net/http"
	"strings"
	"sync"
)

type PageHandler struct {
	PageContentsCache        map[string][]byte
	PageProviders            map[string]PageProvider
	pageContentsCacheRWMutex *sync.RWMutex
	RangeSupported           bool
	CacheSettings            conf.CacheSettingsYaml
}

func NewPageHandler(config conf.ServeYaml) *PageHandler {
	var thePCCMap map[string][]byte
	var theMutex *sync.RWMutex
	if config.CacheSettings.EnableContentsCaching {
		thePCCMap = make(map[string][]byte)
		theMutex = &sync.RWMutex{}
	}
	return &PageHandler{
		PageContentsCache:        thePCCMap,
		PageProviders:            GetProviders(config.CacheSettings.EnableTemplateCaching),
		pageContentsCacheRWMutex: theMutex,
		RangeSupported:           config.RangeSupported,
		CacheSettings:            config.CacheSettings,
	}
}

func (ph *PageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//Provide processing for requests using providers
}

func (ph *PageHandler) PurgeContentsCache(path string, query string) {
	if ph.CacheSettings.EnableContentsCaching {
		if path == "" {
			ph.pageContentsCacheRWMutex.Lock()
			ph.PageContentsCache = make(map[string][]byte)
			ph.pageContentsCacheRWMutex.Unlock()
		} else {
			if strings.HasSuffix(path, "/") {
				ph.pageContentsCacheRWMutex.RLock()
				toDelete := make([]string, len(ph.PageContentsCache))
				theSize := 0
				for cPath := range ph.PageContentsCache {
					dPath := strings.Split(cPath, "?")[0]
					if dPath == path || dPath == path[:len(path)-1] {
						toDelete[theSize] = cPath
						theSize++
					}
				}
				ph.pageContentsCacheRWMutex.RUnlock()
				ph.pageContentsCacheRWMutex.Lock()
				for i := 0; i < theSize; i++ {
					delete(ph.PageContentsCache, toDelete[i])
				}
				ph.pageContentsCacheRWMutex.Unlock()
			} else {
				ph.pageContentsCacheRWMutex.Lock()
				if query == "" {
					delete(ph.PageContentsCache, path)
				} else {
					delete(ph.PageContentsCache, path+"?"+query)
				}
				ph.pageContentsCacheRWMutex.Unlock()
			}
		}
	}
}

func (ph *PageHandler) PurgeTemplateCache(path string) {
	if ph.CacheSettings.EnableTemplateCaching && ph.CacheSettings.EnableTemplateCachePurge {
		if path == "" {
			for _, pageProvider := range ph.PageProviders {
				pageProvider.PurgeTemplate()
			}
		} else {
			if pageProvider, ok := ph.PageProviders[path]; ok {
				pageProvider.PurgeTemplate()
			}
		}
	}
}
