package index

type DataYaml struct {
	About   AboutYaml   `yaml:"about"`
	Entries []EntryYaml `yaml:"entries"`
}
