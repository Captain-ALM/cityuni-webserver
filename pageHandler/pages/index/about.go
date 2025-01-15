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
	r := strings.NewReplacer("#age#", strconv.Itoa(ay.GetAge()), "#birth#", strconv.Itoa(ay.BirthYear), "#year#", strconv.Itoa(time.Now().Year()))
	return template.HTML(r.Replace(ay.Content))
}

func (ay AboutYaml) GetAge() int {
	return time.Now().Year() - ay.BirthYear - 1
}
