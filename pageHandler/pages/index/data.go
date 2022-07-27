package index

import "html/template"

type DataYaml struct {
	HeaderLinks          map[string]template.URL `yaml:"headerLinks"`
	CSSBaseURL           template.URL            `yaml:"cssBaseURL"`
	CSSLightURL          template.URL            `yaml:"cssLightURL"`
	CSSDarkURL           template.URL            `yaml:"cssDarkURL"`
	JScriptURL           template.URL            `yaml:"jScriptURL"`
	NoVideoImageLocation template.URL            `yaml:"noVideoImageLocation"`
	LogoImageLocation    template.URL            `yaml:"logoImageLocation"`
	SunImageLocation     template.URL            `yaml:"sunImageLocation"`
	MoonImageLocation    template.URL            `yaml:"moonImageLocation"`
	About                AboutYaml               `yaml:"about"`
	Entries              []EntryYaml             `yaml:"entries"`
}

func (dy DataYaml) GetHeaderLabels() []string {
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

func (dy DataYaml) GetHeaderLink(headerLabel string) template.URL {
	if dy.HeaderLinks == nil {
		return ""
	}
	return dy.HeaderLinks[headerLabel]
}
