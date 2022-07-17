package index

type DataYaml struct {
	HomeLink      string      `yaml:"homeLink"`
	PortfolioLink string      `yaml:"portfolioLink"`
	CSSBaseURL    string      `yaml:"cssBaseURL"`
	CSSLightURL   string      `yaml:"cssLightURL"`
	CSSDarkURL    string      `yaml:"cssDarkURL"`
	JScriptURL    string      `yaml:"jScriptURL"`
	About         AboutYaml   `yaml:"about"`
	Entries       []EntryYaml `yaml:"entries"`
}
