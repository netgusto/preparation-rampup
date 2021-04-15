package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// List of urls in a file
// Read the file
// loop for each url, call function that calls the url and check HTTP code
// displays the result on stdout (url, call ms, status code, OK or KO)

func main() {

	urls, err := getURLs()
	if err != nil {
		panic(err)
	}

	for round := 0; round < 10; round++ {
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

func getURLs() ([]string, error) {
	data, err := ioutil.ReadFile("/Users/jeromeschneider/Code/rampup/urls.txt")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	nonEmptyLines := []string{}
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return nonEmptyLines, nil
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
