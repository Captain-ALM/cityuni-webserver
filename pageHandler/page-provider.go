package pageHandler

type PageProvider interface {
	GetPath() string
	GetContents(urlParameters map[string]string) (contentType string, contents []byte)
	PurgeTemplate()
}
