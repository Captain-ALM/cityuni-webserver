package index

import (
	"html/template"
	"strings"
)

type DataYaml struct {
	HeaderLinks            map[string]template.URL `yaml:"headerLinks"`
	NavigationEntries      []HeaderEntry           `yaml:"navigationEntries"`
	CSSBaseURL             template.URL            `yaml:"cssBaseURL"`
	CSSHeadersBaseURL      template.URL            `yaml:"cssHeadersBaseURL"`
	CSSMax3HeadersURL      template.URL            `yaml:"cssMax3HeadersURL"`
	CSSOver3HeadersURL     template.URL            `yaml:"cssOver3HeadersURL"`
	CSSLightURL            template.URL            `yaml:"cssLightURL"`
	CSSDarkURL             template.URL            `yaml:"cssDarkURL"`
	JScriptURL             template.URL            `yaml:"jScriptURL"`
	PlayVideoImageLocation template.URL            `yaml:"playVideoImageLocation"`
	NoVideoImageLocation   template.URL            `yaml:"noVideoImageLocation"`
	LogoImageLocation      template.URL            `yaml:"logoImageLocation"`
	SunImageLocation       template.URL            `yaml:"sunImageLocation"`
	MoonImageLocation      template.URL            `yaml:"moonImageLocation"`
	SortImageLocation      template.URL            `yaml:"sortImageLocation"`
	About                  AboutYaml               `yaml:"about"`
	Entries                []EntryYaml             `yaml:"entries"`
}

type HeaderEntry struct {
	Label  string            `yaml:"label"`
	URL    template.URL      `yaml:"url"`
	target template.HTMLAttr `yaml:"target"`
}

func (dy DataYaml) GetAllHeaderLabels() (toReturn []string) {
	if dy.HeaderLinks == nil || len(dy.HeaderLinks) < 1 {
		if dy.NavigationEntries == nil {
			return []string{}
		}
		toReturn = make([]string, len(dy.NavigationEntries))
		for i, entry := range dy.NavigationEntries {
			toReturn[i] = entry.Label
		}
		return
	}
	toReturn = make([]string, len(dy.HeaderLinks))
	i := 0
	for key := range dy.HeaderLinks {
		toReturn[i] = key
		i++
	}
	return
}

func (DataYaml) GetHeaderLabels(labels []string) []string {
	if labels == nil {
		return []string{}
	}
	toReturn := make([]string, min(len(labels), 3))
	for i, key := range labels {
		if i > 2 {
			break
		}
		toReturn[i] = key
	}
	return toReturn
}

func (DataYaml) GetHeaderLabelsExtra(labels []string) []string {
	if len(labels) < 4 {
		return []string{}
	}
	toReturn := make([]string, len(labels)-3)
	for i, key := range labels {
		if i < 3 {
			continue
		}
		toReturn[i-3] = key
	}
	return toReturn
}

func (dy DataYaml) GetHeaderLink(headerLabel string) template.URL {
	if dy.HeaderLinks == nil {
		if dy.NavigationEntries == nil {
			return ""
		}
		for _, entry := range dy.NavigationEntries {
			if entry.Label == headerLabel {
				return entry.URL
			}
		}
		return ""
	}
	return dy.HeaderLinks[headerLabel]
}

func (dy DataYaml) GetHeaderTarget(headerLabel string) template.HTMLAttr {
	if dy.HeaderLinks == nil {
		if dy.NavigationEntries == nil {
			return ""
		}
		for _, entry := range dy.NavigationEntries {
			if entry.Label == headerLabel {
				if strings.EqualFold(string(entry.target), "_self") {
					return ""
				}
				return entry.target
			}
		}
		return ""
	}
	return ""
}

func (dy DataYaml) GetHeaderCount() int {
	if len(dy.HeaderLinks) > 0 {
		return len(dy.HeaderLinks)
	}
	return len(dy.NavigationEntries)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
