package main

import "time"

type Ping struct {
	Moment         time.Time
	Duration       time.Duration
	URL            string
	StatusCode     *int
	TransportError error
	Tries          int
}
