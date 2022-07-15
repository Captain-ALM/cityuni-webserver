package pageHandler

import (
	"golang.captainalm.com/cityuni-webserver/conf"
	"golang.captainalm.com/cityuni-webserver/pageHandler/utils"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PageHandler struct {
	PageContentsCache        map[string]*CachedPage
	PageProviders            map[string]PageProvider
	pageContentsCacheRWMutex *sync.RWMutex
	RangeSupported           bool
	FilterURLQueries         bool
	CacheSettings            conf.CacheSettingsYaml
}

type CachedPage struct {
	Content     []byte
	ContentType string
	LastMod     time.Time
}

func NewPageHandler(config conf.ServeYaml) *PageHandler {
	var thePCCMap map[string]*CachedPage
	var theMutex *sync.RWMutex
	if config.CacheSettings.EnableContentsCaching {
		thePCCMap = make(map[string]*CachedPage)
		theMutex = &sync.RWMutex{}
	}
	return &PageHandler{
		PageContentsCache:        thePCCMap,
		PageProviders:            GetProviders(config.CacheSettings.EnableTemplateCaching, config.DataStorage),
		pageContentsCacheRWMutex: theMutex,
		RangeSupported:           config.RangeSupported,
		FilterURLQueries:         config.FilterURLQueries,
		CacheSettings:            config.CacheSettings,
	}
}

func (ph *PageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	actualPagePath := strings.TrimRight(request.URL.Path, "/")
	queryCollection, actualQueries := ph.GetCleanQuery(request)

	var pageContent []byte
	var pageContentType string
	var lastMod time.Time

	if ph.CacheSettings.EnableContentsCaching {
		cached := ph.getPageFromCache(request.URL, actualQueries)
		if cached != nil {
			pageContent = cached.Content
			pageContentType = cached.ContentType
			lastMod = cached.LastMod
		}
	}

	if pageContentType == "" {
		if provider := ph.PageProviders[actualPagePath]; provider != nil {
			var canCache bool
			pageContentType, pageContent, canCache = provider.GetContents(queryCollection)
			lastMod = provider.GetLastModified()
			if pageContentType != "" && canCache && ph.CacheSettings.EnableContentsCaching {
				ph.setPageToCache(request.URL, actualQueries, &CachedPage{
					Content:     pageContent,
					ContentType: pageContentType,
					LastMod:     lastMod,
				})
			}
		}
	}

	allowedMethods := ph.getAllowedMethodsForPath(request.URL.Path)
	allowed := false
	if request.Method != http.MethodOptions {
		for _, method := range allowedMethods {
			if method == request.Method {
				allowed = true
				break
			}
		}
	}

	if allowed {

		if pageContentType == "" {
			utils.WriteResponseHeaderCanWriteBody(request.Method, writer, http.StatusNotFound, "Page Not Found")
		} else {

			switch request.Method {
			case http.MethodGet, http.MethodHead:

				writer.Header().Set("Content-Type", pageContentType)
				writer.Header().Set("Content-Length", strconv.Itoa(len(pageContent)))
				utils.SetLastModifiedHeader(writer.Header(), lastMod)
				utils.SetCacheHeaderWithAge(writer.Header(), ph.CacheSettings.MaxAge, lastMod)
				theETag := utils.GetValueForETagUsingByteArray(pageContent)
				writer.Header().Set("ETag", theETag)

				if utils.ProcessSupportedPreconditionsForNext(writer, request, lastMod, theETag, ph.CacheSettings.NotModifiedResponseUsingLastModified, ph.CacheSettings.NotModifiedResponseUsingETags) {

					httpRangeParts := utils.ProcessRangePreconditions(int64(len(pageContent)), writer, request, lastMod, theETag, ph.RangeSupported)
					if httpRangeParts != nil {
						if len(httpRangeParts) <= 1 {
							var theWriter io.Writer = writer
							if len(httpRangeParts) == 1 {
								theWriter = utils.NewPartialRangeWriter(theWriter, httpRangeParts[0])
							}
							_, _ = theWriter.Write(pageContent)
						} else {
							multWriter := multipart.NewWriter(writer)
							writer.Header().Set("Content-Type", "multipart/byteranges; boundary="+multWriter.Boundary())
							for _, currentPart := range httpRangeParts {
								mimePart, err := multWriter.CreatePart(textproto.MIMEHeader{
									"Content-Range": {currentPart.ToField(int64(len(pageContent)))},
									"Content-Type":  {"text/plain; charset=utf-8"},
								})
								if err != nil {
									break
								}
								_, err = mimePart.Write(pageContent[currentPart.Start : currentPart.Start+currentPart.Length])
								if err != nil {
									break
								}
							}
							_ = multWriter.Close()
						}
					}
				}
			case http.MethodDelete:
				ph.PurgeTemplateCache(actualPagePath)
				ph.PurgeContentsCache(request.URL.Path, actualQueries)
				utils.SetNeverCacheHeader(writer.Header())
				utils.WriteResponseHeaderCanWriteBody(request.Method, writer, http.StatusOK, "")
			}
		}
	} else {

		theAllowHeaderContents := ""
		for _, method := range allowedMethods {
			theAllowHeaderContents += method + ", "
		}

		writer.Header().Set("Allow", strings.TrimSuffix(theAllowHeaderContents, ", "))
		if request.Method == http.MethodOptions {
			utils.WriteResponseHeaderCanWriteBody(request.Method, writer, http.StatusOK, "")
		} else {
			utils.WriteResponseHeaderCanWriteBody(request.Method, writer, http.StatusMethodNotAllowed, "")
		}
	}
}

func (ph *PageHandler) GetCleanQuery(request *http.Request) (url.Values, string) {
	toClean := request.URL.Query()
	provider := ph.PageProviders[request.URL.Path]
	if provider == nil {
		return make(url.Values), ""
	}
	supportedKeys := provider.GetSupportedURLParameters()
	var toDelete []string
	if ph.FilterURLQueries {
		toDelete = make([]string, len(toClean))
	}
	theSize := 0
	theQuery := ""
	for s, v := range toClean {
		noExist := true
		for _, key := range supportedKeys {
			if s == key {
				noExist = false
				break
			}
		}
		if noExist {
			if ph.FilterURLQueries {
				toDelete[theSize] = s
				theSize++
			}
		} else {
			for _, i := range v {
				if i == "" {
					theQuery += s + "&"
				} else {
					theQuery += s + "=" + i + "&"
				}
			}
		}
	}
	if ph.FilterURLQueries {
		for i := 0; i < theSize; i++ {
			delete(toClean, toDelete[i])
		}
	}
	return toClean, strings.TrimRight(theQuery, "&")
}

func (ph *PageHandler) PurgeContentsCache(path string, query string) {
	if ph.CacheSettings.EnableContentsCaching && ph.CacheSettings.EnableContentsCachePurge {
		if path == "" {
			ph.pageContentsCacheRWMutex.Lock()
			ph.PageContentsCache = make(map[string]*CachedPage)
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
func (ph *PageHandler) getPageFromCache(urlIn *url.URL, cleanedQueries string) *CachedPage {
	ph.pageContentsCacheRWMutex.RLock()
	defer ph.pageContentsCacheRWMutex.RUnlock()
	if strings.HasSuffix(urlIn.Path, "/") {
		return ph.PageContentsCache[strings.TrimRight(urlIn.Path, "/")]
	} else {
		if cleanedQueries == "" {
			return ph.PageContentsCache[urlIn.Path]
		} else {
			return ph.PageContentsCache[urlIn.Path+"?"+cleanedQueries]
		}
	}
}

func (ph *PageHandler) setPageToCache(urlIn *url.URL, cleanedQueries string, newPage *CachedPage) {
	ph.pageContentsCacheRWMutex.Lock()
	defer ph.pageContentsCacheRWMutex.Unlock()
	if strings.HasSuffix(urlIn.Path, "/") {
		ph.PageContentsCache[strings.TrimRight(urlIn.Path, "/")] = newPage
	} else {
		if cleanedQueries == "" {
			ph.PageContentsCache[urlIn.Path] = newPage
		} else {
			ph.PageContentsCache[urlIn.Path+"?"+cleanedQueries] = newPage
		}
	}
}

func (ph *PageHandler) getAllowedMethodsForPath(pathIn string) []string {
	if strings.HasSuffix(pathIn, "/") {
		if (ph.CacheSettings.EnableTemplateCaching && ph.CacheSettings.EnableTemplateCachePurge) ||
			(ph.CacheSettings.EnableContentsCaching && ph.CacheSettings.EnableContentsCachePurge) {
			return []string{http.MethodHead, http.MethodGet, http.MethodOptions, http.MethodDelete}
		} else {
			return []string{http.MethodHead, http.MethodGet, http.MethodOptions}
		}
	} else {
		if ph.CacheSettings.EnableContentsCaching && ph.CacheSettings.EnableContentsCachePurge {
			return []string{http.MethodHead, http.MethodGet, http.MethodOptions, http.MethodDelete}
		} else {
			return []string{http.MethodHead, http.MethodGet, http.MethodOptions}
		}
	}
}
