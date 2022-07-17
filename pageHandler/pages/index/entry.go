package index

import (
	"html/template"
	"time"
)

const dateFormat = "2006-01-02"

type EntryYaml struct {
	Name               string    `yaml:"name"`
	Content            string    `yaml:"content"`
	StartDate          time.Time `yaml:"startDate"`
	EndDate            time.Time `yaml:"endDate"`
	VideoLocation      string    `yaml:"videoLocation"`
	VideoContentType   string    `yaml:"videoContentType"`
	ThumbnailLocations []string  `yaml:"thumbnailLocations"`
	ImageLocations     []string  `yaml:"imageLocations"`
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
		return ey.EndDate
	}
}

func (ey EntryYaml) GetContent() template.HTML {
	return template.HTML(ey.Content)
}

func (ey EntryYaml) GetDuration() time.Duration {
	return ey.GetEndTime().Sub(ey.StartDate).Truncate(time.Second)
}
