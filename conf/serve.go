package conf

type ServeYaml struct {
	DataStorage    string            `yaml:"dataStorage"`
	Domains        []string          `yaml:"domains"`
	RangeSupported bool              `yaml:"rangeSupported"`
	CacheSettings  CacheSettingsYaml `yaml:"cacheSettings"`
}
