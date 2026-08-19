package network

import (
	"crypto/tls"
	"net/http"
)

// @inject
func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
