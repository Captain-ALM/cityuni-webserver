package pageHandler

import "net/url"

type PageProvider interface {
	GetPath() string
	GetSupportedURLParameters() []string
	GetContents(urlParameters url.Values) (contentType string, contents []byte)
	PurgeTemplate()
}
