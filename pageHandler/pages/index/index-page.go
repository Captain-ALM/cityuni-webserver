package index

import (
	"golang.captainalm.com/cityuni-webserver/utils/io"
	"gopkg.in/yaml.v3"
	"html/template"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const templateName = "index.go.html"
const yamlName = "index.go.yml"

func NewPage(dataStore string, cacheTemplates bool) *Page {
	var ptm *sync.Mutex
	var sdm *sync.Mutex
	if cacheTemplates {
		ptm = &sync.Mutex{}
		sdm = &sync.Mutex{}
	}
	pageToReturn := &Page{
		DataStore:         dataStore,
		StoredDataMutex:   sdm,
		PageTemplateMutex: ptm,
	}
	return pageToReturn
}

type Page struct {
	DataStore            string
	StoredDataMutex      *sync.Mutex
	StoredData           *DataYaml
	LastModifiedData     time.Time
	PageTemplateMutex    *sync.Mutex
	PageTemplate         *template.Template
	LastModifiedTemplate time.Time
}

func (p *Page) GetPath() string {
	return "/index.go"
}

func (p *Page) GetLastModified() time.Time {
	if p.LastModifiedData.After(p.LastModifiedTemplate) {
		return p.LastModifiedData
	} else {
		return p.LastModifiedTemplate
	}
}

func (p *Page) GetCacheIDExtension(urlParameters url.Values) string {
	toReturn := ""
	if urlParameters.Has("order") {
		if theParameter := strings.ToLower(urlParameters.Get("order")); theParameter == "start" || theParameter == "end" || theParameter == "name" || theParameter == "duration" {
			toReturn += "order=" + theParameter + "&"
		}
	}
	if urlParameters.Has("sort") {
		if theParameter := strings.ToLower(urlParameters.Get("sort")); theParameter == "asc" || theParameter == "ascending" || theParameter == "desc" || theParameter == "descending" {
			toReturn += "sort=" + theParameter + "&"
		}
	}
	if urlParameters.Has("light") {
		toReturn += "light"
	}
	return strings.TrimRight(toReturn, "&")
}

func (p *Page) GetContents(urlParameters url.Values) (contentType string, contents []byte, canCache bool) {
	theTemplate, err := p.getPageTemplate()
	if err != nil {
		return "text/plain", []byte("Cannot Get Index.\r\n" + err.Error()), false
	}
	theData, err := p.getPageData()
	if err != nil {
		return "text/plain", []byte("Cannot Get Data.\r\n" + err.Error()), false
	}
	theMarshal := &Marshal{
		Data:  *theData,
		Light: urlParameters.Has("light"),
	}
	switch strings.ToLower(urlParameters.Get("order")) {
	case "end":
		theMarshal.OrderEndDate = getSortValue(strings.ToLower(urlParameters.Get("sort")))
	case "name":
		theMarshal.OrderName = getSortValue(strings.ToLower(urlParameters.Get("sort")))
	case "duration":
		theMarshal.OrderDuration = getSortValue(strings.ToLower(urlParameters.Get("sort")))
	default:
		theMarshal.OrderStartDate = getSortValue(strings.ToLower(urlParameters.Get("sort")))
	}
	theBuffer := &io.BufferedWriter{}
	err = theTemplate.ExecuteTemplate(theBuffer, templateName, theMarshal)
	if err != nil {
		return "text/plain", []byte("Cannot Get Page.\r\n" + err.Error()), false
	}
	return "text/html", theBuffer.Data, true
}

func (p *Page) PurgeTemplate() {
	if p.PageTemplateMutex != nil {
		p.PageTemplateMutex.Lock()
		p.PageTemplate = nil
		p.PageTemplateMutex.Unlock()
	}
	if p.StoredDataMutex != nil {
		p.StoredDataMutex.Lock()
		p.StoredData = nil
		p.StoredDataMutex.Unlock()
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
		stat, err := os.Stat(thePath)
		if err != nil {
			return nil, err
		}
		p.LastModifiedTemplate = stat.ModTime()
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

func (p *Page) getPageData() (*DataYaml, error) {
	if p.StoredDataMutex != nil {
		p.StoredDataMutex.Lock()
		defer p.StoredDataMutex.Unlock()
	}
	if p.StoredData == nil {
		thePath := yamlName
		if p.DataStore != "" {
			thePath = path.Join(p.DataStore, thePath)
		}
		stat, err := os.Stat(thePath)
		if err != nil {
			return nil, err
		}
		p.LastModifiedData = stat.ModTime()
		fileHandle, err := os.Open(thePath)
		if err != nil {
			return nil, err
		}
		dataYaml := &DataYaml{}
		decoder := yaml.NewDecoder(fileHandle)
		err = decoder.Decode(dataYaml)
		if err != nil {
			_ = fileHandle.Close()
			return nil, err
		}
		err = fileHandle.Close()
		if err != nil {
			return nil, err
		}
		if p.StoredDataMutex != nil {
			p.StoredData = dataYaml
		}
		return dataYaml, nil
	} else {
		return p.StoredData, nil
	}
}

func getSortValue(toCheckIn string) int8 {
	if toCheckIn == "desc" || toCheckIn == "descending" {
		return -1
	} else {
		return 1
	}
}
