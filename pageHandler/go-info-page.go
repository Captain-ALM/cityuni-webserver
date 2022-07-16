package pageHandler

import (
	"golang.captainalm.com/cityuni-webserver/conf"
	"golang.captainalm.com/cityuni-webserver/utils/info"
	"golang.captainalm.com/cityuni-webserver/utils/io"
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

func (gipg *goInfoPage) GetCacheIDExtension(urlParameters url.Values) string {
	if urlParameters.Has("full") {
		return "full"
	} else {
		return ""
	}
}

type goInfoTemplateMarshal struct {
	FullOutput         bool
	RegisteredPages    []string
	CachedPages        []string
	ProcessID          int
	ProductLocation    string
	ProductName        string
	ProductDescription string
	BuildVersion       string
	BuildDate          string
	WorkingDirectory   string
	Hostname           string
	PageSize           int
	GoVersion          string
	GoRoutineNum       int
	GoCGoCallNum       int64
	NumCPU             int
	GoRoot             string
	GoMaxProcs         int
	ListenSettings     conf.ListenYaml
	ServeSettings      conf.ServeYaml
	Environment        []string
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
	theBuffer := &io.BufferedWriter{}
	var regPages []string
	var cacPages []string
	env := make([]string, 0)
	if urlParameters.Has("full") {
		regPages = gipg.Handler.GetRegisteredPages()
		cacPages = gipg.Handler.GetCachedPages()
		env = os.Environ()
	} else {
		regPages = make([]string, len(gipg.Handler.PageProviders))
		cacPages = make([]string, gipg.Handler.GetNumberOfCachedPages())
	}
	err = theTemplate.ExecuteTemplate(theBuffer, templateName, &goInfoTemplateMarshal{
		FullOutput:         urlParameters.Has("full"),
		RegisteredPages:    regPages,
		CachedPages:        cacPages,
		ProcessID:          os.Getpid(),
		ProductLocation:    getStringOrError(os.Executable),
		ProductName:        info.BuildName,
		ProductDescription: info.BuildDescription,
		BuildVersion:       info.BuildVersion,
		BuildDate:          info.BuildDate,
		WorkingDirectory:   getStringOrError(os.Getwd),
		Hostname:           getStringOrError(os.Hostname),
		PageSize:           os.Getpagesize(),
		GoVersion:          runtime.Version(),
		GoRoutineNum:       runtime.NumGoroutine(),
		GoCGoCallNum:       runtime.NumCgoCall(),
		NumCPU:             runtime.NumCPU(),
		GoRoot:             runtime.GOROOT(),
		GoMaxProcs:         runtime.GOMAXPROCS(0),
		ListenSettings:     info.ListenSettings,
		ServeSettings:      info.ServeSettings,
		Environment:        env,
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

func getStringOrError(funcIn func() (string, error)) string {
	toReturn, err := funcIn()
	if err == nil {
		return toReturn
	} else {
		return "Error: " + err.Error()
	}
}
