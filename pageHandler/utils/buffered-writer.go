package utils

import (
	"crypto"
	"encoding/hex"
)

type BufferedWriter struct {
	Data []byte
}

func (c *BufferedWriter) Write(p []byte) (n int, err error) {
	c.Data = append(c.Data, p...)
	return len(p), nil
}

func (c *BufferedWriter) GetHashString() string {
	theHash := crypto.SHA1.New()
	_, _ = theHash.Write(c.Data)
	theSum := theHash.Sum(nil)
	theHash.Reset()
	return hex.EncodeToString(theSum)
}
