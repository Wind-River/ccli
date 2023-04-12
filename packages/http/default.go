package http

import (
	"crypto/tls"
	"net/http"
)

var DefaultClient *http.Client

func init() {
	DefaultClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
