package api

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	Version  = `5.120`
	Endpoint = `https://api.vk.com/method`
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func Rnd() int64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Int63()
}
