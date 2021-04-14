package main

import (
	"net/http"
	"time"
)

type URLGetter interface {
	GetURL(url string) (*int, time.Duration, error)
}

type URLGetterReal struct{}

func (ug URLGetterReal) GetURL(url string) (*int, time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return nil, duration, err
	}

	return &resp.StatusCode, duration, nil
}
