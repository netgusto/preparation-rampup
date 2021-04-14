package main

import (
	"fmt"
	"net/http"
	"time"
)

// List of urls in a file (in an array fine too)
// Read the file
// loop for each url, call function that calls the url and check HTTP code
// displays the result on stdout (url, call ms, status code, OK or KO)
// Write a synchronous HTTP probe the simplest way possible

func main() {

	urls := getURLs()

	for round := 1; round <= 10; round++ {
		println("# ROUND ", round)
		for _, url := range urls {
			code, duration, err := probe(url)

			var outcome string
			if err == nil {
				outcome = "SUCCESS"
			} else {
				outcome = "ERROR, " + err.Error()
			}

			fmt.Printf("%s; STATUS=%d; DURATION=%s; OUTCOME=%s\n", url, code, duration, outcome)
		}
	}
}

func getURLs() []string {
	return []string{
		"https://www.algolia.com",
		"https://d85-usw-1.algolia.net/1/isalive",
		"https://d85-usw-2.algolia.net/1/isalive",
		"https://d85-usw-3.algolia.net/1/isalive",
	}
}

func probe(url string) (int, time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return 0, duration, err
	}

	return resp.StatusCode, duration, nil
}
