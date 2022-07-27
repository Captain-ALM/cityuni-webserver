package index

import (
	"html/template"
	"strconv"
	"strings"
	"time"
)

type AboutYaml struct {
	Title             string       `yaml:"title"`
	Content           string       `yaml:"content"`
	ThumbnailLocation template.URL `yaml:"thumbnailLocation"`
	ImageLocation     template.URL `yaml:"imageLocation"`
	ImageAltText      string       `yaml:"imageAltText"`
	BirthYear         int          `yaml:"birthYear"`
	ContactEmail      string       `yaml:"contactEmail"`
}

func (ay AboutYaml) GetContent() template.HTML {
	return template.HTML(strings.ReplaceAll(strings.ReplaceAll(ay.Content, "#age#", strconv.Itoa(ay.GetAge())), "#birth#", strconv.Itoa(ay.BirthYear)))
}

func (ay AboutYaml) GetAge() int {
	return time.Now().Year() - ay.BirthYear - 1
}
