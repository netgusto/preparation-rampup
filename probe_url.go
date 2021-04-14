package main

import (
	"net/http"
	"time"
)

func probeURL(url string, maxTries int, urlGetter URLGetter) Ping {

	start := time.Now().Round(0)

	var code *int
	var duration time.Duration
	var err error
	var tries int

	for k := 0; k < maxTries; k++ {
		tries++
		code, duration, err = urlGetter.GetURL(url)
		if err == nil && code != nil && *code == http.StatusOK {
			// we got one successful call
			break
		}
	}

	return Ping{
		Moment:         start,
		Duration:       duration,
		URL:            url,
		StatusCode:     code,
		TransportError: err,
		Tries:          tries,
	}
}
