package pageHandler

import (
	"golang.captainalm.com/cityuni-webserver/conf"
	"golang.captainalm.com/cityuni-webserver/pageHandler/utils"
	"golang.captainalm.com/cityuni-webserver/utils/info"
	"html/template"
	"net/url"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const templateName = "goinfo.go.html"

func newGoInfoPage(handlerIn *PageHandler, dataStore string, cacheTemplates bool) *goInfoPage {
	var ptm *sync.Mutex
	if cacheTemplates {
		ptm = &sync.Mutex{}
	}
	return &goInfoPage{
		Handler:           handlerIn,
		DataStore:         dataStore,
		CacheTemplate:     cacheTemplates,
		PageTemplateMutex: ptm,
	}
}

type goInfoPage struct {
	Handler           *PageHandler
	DataStore         string
	CacheTemplate     bool
	PageTemplateMutex *sync.Mutex
	PageTemplate      *template.Template
}

type goInfoTemplateMarshal struct {
	FullOutput         bool
	RegisteredPages    []string
	CachedPages        []string
	ProductName        string
	ProductDescription string
	BuildVersion       string
	BuildDate          string
	GoVersion          string
	GoRoutineNum       int
	GoCGoCallNum       int64
	NumCPU             int
	GoRoot             string
	GoMaxProcs         int
	ListenSettings     conf.ListenYaml
	ServeSettings      conf.ServeYaml
}

func (gipg *goInfoPage) GetPath() string {
	return "/goinfo.go"
}

func (gipg *goInfoPage) GetSupportedURLParameters() []string {
	return []string{"full"}
}

func (gipg *goInfoPage) GetLastModified() time.Time {
	return time.Now()
}

func (gipg *goInfoPage) GetContents(urlParameters url.Values) (contentType string, contents []byte, canCache bool) {
	theTemplate, err := gipg.getPageTemplate()
	if err != nil {
		return "text/plain", []byte("Cannot Get Info.\r\n" + err.Error()), false
	}
	theBuffer := &utils.BufferedWriter{}
	err = theTemplate.ExecuteTemplate(theBuffer, templateName, &goInfoTemplateMarshal{
		FullOutput:         urlParameters.Has("full"),
		RegisteredPages:    gipg.Handler.GetRegisteredPages(),
		CachedPages:        gipg.Handler.GetCachedPages(),
		ProductName:        info.BuildName,
		ProductDescription: info.BuildDescription,
		BuildVersion:       info.BuildVersion,
		BuildDate:          info.BuildDate,
		GoVersion:          runtime.Version(),
		GoRoutineNum:       runtime.NumGoroutine(),
		GoCGoCallNum:       runtime.NumCgoCall(),
		NumCPU:             runtime.NumCPU(),
		GoRoot:             runtime.GOROOT(),
		GoMaxProcs:         runtime.GOMAXPROCS(0),
		ListenSettings:     info.ListenSettings,
		ServeSettings:      info.ServeSettings,
	})
	if err != nil {
		return "text/plain", []byte("Cannot Get Info.\r\n" + err.Error()), false
	}
	return "text/html", theBuffer.Data, false
}

func (gipg *goInfoPage) PurgeTemplate() {
	if gipg.CacheTemplate {
		gipg.PageTemplateMutex.Lock()
		gipg.PageTemplate = nil
		gipg.PageTemplateMutex.Unlock()
	}
}

func (gipg *goInfoPage) getPageTemplate() (*template.Template, error) {
	if gipg.CacheTemplate {
		gipg.PageTemplateMutex.Lock()
		defer gipg.PageTemplateMutex.Unlock()
	}
	if gipg.PageTemplate == nil {
		thePath := templateName
		if gipg.DataStore != "" {
			thePath = path.Join(gipg.DataStore, thePath)
		}
		loadedData, err := os.ReadFile(thePath)
		if err != nil {
			return nil, err
		}
		tmpl, err := template.New(templateName).Parse(string(loadedData))
		if err != nil {
			return nil, err
		}
		if gipg.CacheTemplate {
			gipg.PageTemplate = tmpl
		}
		return tmpl, nil
	} else {
		return gipg.PageTemplate, nil
	}
}
