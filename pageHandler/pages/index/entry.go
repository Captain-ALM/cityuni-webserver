package index

import (
	"golang.captainalm.com/cityuni-webserver/utils/yaml"
	"html/template"
	"time"
)

const dateFormat = "2006-01"

type EntryYaml struct {
	Name               string        `yaml:"name"`
	Content            string        `yaml:"content"`
	StartDate          yaml.DateType `yaml:"startDate"`
	EndDate            yaml.DateType `yaml:"endDate"`
	VideoLocation      string        `yaml:"videoLocation"`
	VideoContentType   string        `yaml:"videoContentType"`
	ThumbnailLocations []string      `yaml:"thumbnailLocations"`
	ImageLocations     []string      `yaml:"imageLocations"`
}

func (ey EntryYaml) GetStartDate() string {
	return ey.StartDate.Format(dateFormat)
}

func (ey EntryYaml) GetEndDate() string {
	if ey.EndDate.IsZero() {
		return ""
	} else {
		return ey.EndDate.Format(dateFormat)
	}
}

func (ey EntryYaml) GetEndTime() time.Time {
	if ey.EndDate.IsZero() {
		return time.Now()
	} else {
		return ey.EndDate.Time
	}
}

func (ey EntryYaml) GetContent() template.HTML {
	return template.HTML(ey.Content)
}

func (ey EntryYaml) GetDuration() time.Duration {
	return ey.GetEndTime().Sub(ey.StartDate.Time).Truncate(time.Second)
}
