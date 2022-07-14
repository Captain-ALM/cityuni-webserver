package conf

type ServeYaml struct {
	Domains        []string          `yaml:"domains"`
	RangeSupported bool              `yaml:"rangeSupported"`
	CacheSettings  CacheSettingsYaml `yaml:"cacheSettings"`
}
