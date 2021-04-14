package main

import (
	"fmt"
	"io/ioutil"
)

// List of urls in a file
// Read the file
// loop for each url, call function that calls the url and check HTTP code
// displays the result on stdout (url, call ms, status code, OK or KO)

const nbTriesPerURL = 3

func main() {
	data, err := ioutil.ReadFile("/Users/jeromeschneider/Code/rampup/urls.txt")
	if err != nil {
		panic(err)
	}

	urls, err := getURLList(data)
	if err != nil {
		panic(err)
	}

	urlGetter := URLGetterReal{}

	for round := 1; round <= 10; round++ {
		println("# ROUND ", round)
		for _, url := range urls {
			ping := probeURL(url, nbTriesPerURL, urlGetter)
			fmt.Printf("%+v\n", ping)
		}
	}
}
