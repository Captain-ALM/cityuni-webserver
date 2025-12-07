package index

import (
	"html/template"
	"strconv"
	"strings"
	"time"
)

type AboutYaml struct {
	Title             string                  `yaml:"title"`
	Content           string                  `yaml:"content"`
	ThumbnailLocation template.URL            `yaml:"thumbnailLocation"`
	ImageLocation     template.URL            `yaml:"imageLocation"`
	ImageAltText      string                  `yaml:"imageAltText"`
	BirthYear         int                     `yaml:"birthYear"`
	ContactEmails     []ContactEmailYaml      `yaml:"contactEmails"`
	ExtraLinks        map[string]template.URL `yaml:"extraLinks"`
}

type ContactEmailYaml struct {
	Email      string       `yaml:"email"`
	GPGUrl     template.URL `yaml:"gpgURL"`
	Thumbprint string       `yaml:"thumbprint"`
}

func (ay AboutYaml) GetContent() template.HTML {
	r := strings.NewReplacer("#age#", strconv.Itoa(ay.GetAge()), "#birth#", strconv.Itoa(ay.BirthYear), "#year#", strconv.Itoa(time.Now().Year()))
	return template.HTML(r.Replace(ay.Content))
}

func (ay AboutYaml) GetAge() int {
	return time.Now().Year() - ay.BirthYear - 1
}

func (ay AboutYaml) GetContactEmail() ContactEmailYaml {
	if len(ay.ContactEmails) < 1 {
		return ContactEmailYaml{}
	}
	return ay.ContactEmails[0]
}

func (ay AboutYaml) GetAlternateContactEmail() ContactEmailYaml {
	if len(ay.ContactEmails) < 2 {
		return ContactEmailYaml{}
	}
	return ay.ContactEmails[1]
}

func (ay AboutYaml) GetExtraLinkLabels() []string {
	if ay.ExtraLinks == nil {
		return []string{}
	}
	toReturn := make([]string, len(ay.ExtraLinks))
	i := 0
	for key := range ay.ExtraLinks {
		toReturn[i] = key
		i++
	}
	return toReturn
}

func (ay AboutYaml) GetExtraLinkLink(linkLabel string) template.URL {
	if ay.ExtraLinks == nil {
		return ""
	}
	return ay.ExtraLinks[linkLabel]
}

func (ay AboutYaml) GetExtraLinkCount() int {
	return len(ay.ExtraLinks)
}
