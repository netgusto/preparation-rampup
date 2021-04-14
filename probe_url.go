package main

import (
	"net/http"
	"time"
)

func probeURL(url string, maxTries int) Ping {

	start := time.Now().Round(0)

	var code *int
	var duration time.Duration
	var err error
	var tries int

	for k := 0; k < maxTries; k++ {
		tries++
		code, duration, err = measureURLCall(url)
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

func measureURLCall(url string) (*int, time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return nil, duration, err
	}

	return &resp.StatusCode, duration, nil
}
