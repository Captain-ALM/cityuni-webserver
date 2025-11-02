package index

import (
	"crypto/sha1"
	"encoding/hex"
	"golang.captainalm.com/cityuni-webserver/utils/yaml"
	"html/template"
	"math"
	"net/http"
	"strings"
	"time"
)

const dateFormat = "01/2006"

var hash = sha1.New

type EntryYaml struct {
	Name                   string         `yaml:"name"`
	Content                string         `yaml:"content"`
	StartDate              yaml.DateType  `yaml:"startDate"`
	EndDate                yaml.DateType  `yaml:"endDate"`
	VideoLocation          template.URL   `yaml:"videoLocation"`
	VideoContentType       string         `yaml:"videoContentType"`
	ThumbnailLocations     []template.URL `yaml:"thumbnailLocations"`
	ImageLocations         []template.URL `yaml:"imageLocations"`
	ImageAltTexts          []string       `yaml:"imageAltTexts"`
	VideoThumbnailLocation template.URL   `yaml:"videoThumbnailLocation"`
}

type ImageReference struct {
	ThumbnailLocation template.URL
	ImageLocation     template.URL
	ImageAltText      string
}

func (ey EntryYaml) GetEntryAnchor() string {
	h := hash()
	h.Write([]byte(ey.Name))
	return "entry_" + hex.EncodeToString(h.Sum(nil))
}

func (ey EntryYaml) GetVideoThumbnail(usual template.URL) template.URL {
	if ey.VideoThumbnailLocation == "" {
		return usual
	} else {
		return ey.VideoThumbnailLocation
	}
}

func (ey EntryYaml) IsVideoLink() bool {
	return strings.EqualFold(ey.VideoContentType, "text/uri-list")
}

func (ey EntryYaml) GetStartDate() string {
	return ey.StartDate.Format(dateFormat)
}

func (ey EntryYaml) GetStartDateHTML() string {
	return ey.StartDate.Format(http.TimeFormat)
}

func (ey EntryYaml) GetEndDate() string {
	if ey.EndDate.IsZero() {
		return ""
	} else {
		return ey.EndDate.Format(dateFormat)
	}
}

func (ey EntryYaml) GetEndDateHTML() string {
	return ey.GetEndTime().Format(http.TimeFormat)
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

func (ey EntryYaml) GetInt64Duration() int64 {
	return int64(ey.GetDuration())
}

func (ey EntryYaml) GetImageCount() int {
	return int(math.Min(math.Min(float64(len(ey.ThumbnailLocations)), float64(len(ey.ImageLocations))), float64(len(ey.ImageAltTexts))))
}

func (ey EntryYaml) GetImages() []ImageReference {
	toReturn := make([]ImageReference, ey.GetImageCount())
	for i := 0; i < len(ey.ThumbnailLocations) && i < len(ey.ImageLocations) && i < len(ey.ImageAltTexts); i++ {
		toReturn[i] = ImageReference{
			ThumbnailLocation: ey.ThumbnailLocations[i],
			ImageLocation:     ey.ImageLocations[i],
			ImageAltText:      ey.ImageAltTexts[i],
		}
	}
	return toReturn
}
