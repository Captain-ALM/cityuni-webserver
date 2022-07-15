package pageHandler

import (
	"net/url"
	"time"
)

type PageProvider interface {
	GetPath() string
	GetSupportedURLParameters() []string
	GetLastModified() time.Time
	GetContents(urlParameters url.Values) (contentType string, contents []byte, canCache bool)
	PurgeTemplate()
}
