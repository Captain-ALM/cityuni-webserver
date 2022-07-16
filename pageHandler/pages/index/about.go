package index

import (
	"html/template"
	"time"
)

type AboutYaml struct {
	Title             string `yaml:"title"`
	Content           string `yaml:"content"`
	ThumbnailLocation string `yaml:"thumbnailLocation"`
	ImageLocation     string `yaml:"imageLocation"`
	BirthYear         int    `yaml:"birthYear"`
	ContactEmail      string `yaml:"contactEmail"`
}

func (ay AboutYaml) GetContent() template.HTML {
	return template.HTML(ay.Content)
}

func (ay AboutYaml) GetAge() int {
	return time.Now().Year() - ay.BirthYear - 1
}
