package index

import (
	"html/template"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const templateName = "index.go.html"

func NewPage(dataStore string, cacheTemplates bool) *Page {
	var ptm *sync.Mutex
	if cacheTemplates {
		ptm = &sync.Mutex{}
	}
	pageToReturn := &Page{
		DataStore:         dataStore,
		PageTemplateMutex: ptm,
	}
	if !cacheTemplates {
		_, _ = pageToReturn.getPageTemplate()
	}
	return pageToReturn
}

type Page struct {
	DataStore         string
	PageTemplateMutex *sync.Mutex
	PageTemplate      *template.Template
	LastModified      time.Time
}

func (p *Page) GetPath() string {
	return "/index.go"
}

func (p *Page) GetLastModified() time.Time {
	return p.LastModified
}

func (p *Page) GetCacheIDExtension(urlParameters url.Values) string {
	toReturn := ""
	if urlParameters.Has("order") {
		if theParameter := strings.ToLower(urlParameters.Get("order"));
			theParameter == "start" || theParameter == "end" || theParameter == "name" || theParameter == "duration" {
			toReturn += "order=" + theParameter + "&"
		}
	}
	if urlParameters.Has("sort") {
		if theParameter := strings.ToLower(urlParameters.Get("sort"));
			theParameter == "asc" || theParameter == "ascending" || theParameter == "desc" || theParameter == "descending" {
			toReturn += "sort=" + theParameter
		}
	}
	return strings.TrimRight(toReturn, "&")
}

func (p *Page) GetContents(urlParameters url.Values) (contentType string, contents []byte, canCache bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Page) PurgeTemplate() {
	if p.PageTemplateMutex != nil {
		p.PageTemplateMutex.Lock()
		p.PageTemplate = nil
		p.PageTemplateMutex.Unlock()
	}
}

func (p *Page) getPageTemplate() (*template.Template, error) {
	if p.PageTemplateMutex != nil {
		p.PageTemplateMutex.Lock()
		defer p.PageTemplateMutex.Unlock()
	}
	if p.PageTemplate == nil {
		thePath := templateName
		if p.DataStore != "" {
			thePath = path.Join(p.DataStore, thePath)
		}
		loadedData, err := os.ReadFile(thePath)
		if err != nil {
			return nil, err
		}
		tmpl, err := template.New(templateName).Parse(string(loadedData))
		if err != nil {
			return nil, err
		}
		if p.PageTemplateMutex != nil {
			p.PageTemplate = tmpl
		}
		return tmpl, nil
	} else {
		return p.PageTemplate, nil
	}
}
