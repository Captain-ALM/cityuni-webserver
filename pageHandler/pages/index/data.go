package index

type DataYaml struct {
	HeaderLinks map[string]string `yaml:"headerLinks"`
	CSSBaseURL  string            `yaml:"cssBaseURL"`
	CSSLightURL string            `yaml:"cssLightURL"`
	CSSDarkURL  string            `yaml:"cssDarkURL"`
	JScriptURL  string            `yaml:"jScriptURL"`
	About       AboutYaml         `yaml:"about"`
	Entries     []EntryYaml       `yaml:"entries"`
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

func (dy DataYaml) GetHeaderLink(headerLabel string) string {
	if dy.HeaderLinks == nil {
		return ""
	}
	return dy.HeaderLinks[headerLabel]
}
