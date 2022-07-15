package pages

import (
	"net/url"
	"time"
)

var startTime = time.Now()

func NewTestPage() *TestPage {
	return &TestPage{}
}

type TestPage struct {
}

func (tp *TestPage) GetPath() string {
	return "/test.go"
}

func (tp *TestPage) GetSupportedURLParameters() []string {
	return []string{"test"}
}

func (tp *TestPage) GetLastModified() time.Time {
	return startTime
}

func (tp *TestPage) GetContents(urlParameters url.Values) (contentType string, contents []byte) {
	if val, ok := urlParameters["test"]; ok {
		if len(val) > 0 {
			return "text/plain", ([]byte)("Testing!\r\n" + val[0])
		}
	}
	return "text/plain", ([]byte)("Testing!")
}

func (tp *TestPage) PurgeTemplate() {
}
