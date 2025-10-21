package index

import "html/template"

type DataYaml struct {
	HeaderLinks            map[string]template.URL `yaml:"headerLinks"`
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

func (dy DataYaml) GetAllHeaderLabels() []string {
	if dy.HeaderLinks == nil {
		return []string{}
	}
	toReturn := make([]string, len(dy.HeaderLinks))
	i := 0
	for key := range dy.HeaderLinks {
		toReturn[i] = key
		i++
	}
	return toReturn
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
		return ""
	}
	return dy.HeaderLinks[headerLabel]
}

func (dy DataYaml) GetHeaderCount() int {
	return len(dy.HeaderLinks)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
