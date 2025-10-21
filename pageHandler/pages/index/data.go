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

func (dy DataYaml) GetHeaderLabels() []string {
	if dy.HeaderLinks == nil {
		return []string{}
	}
	toReturn := make([]string, min(len(dy.HeaderLinks), 3))
	i := 0
	for key := range dy.HeaderLinks {
		toReturn[i] = key
		i++
		if i > 2 {
			break
		}
	}
	return toReturn
}

func (dy DataYaml) GetHeaderLabelsExtra() []string {
	if len(dy.HeaderLinks) < 4 {
		return []string{}
	}
	toReturn := make([]string, len(dy.HeaderLinks)-3)
	i := 0
	for key := range dy.HeaderLinks {
		if i < 3 {
			i++
			continue
		}
		toReturn[i-3] = key
		i++
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
