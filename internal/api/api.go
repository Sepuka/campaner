package api

import "net/http"

const (
	Version  = `5.120`
	Endpoint = `https://api.vk.com/method`
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
