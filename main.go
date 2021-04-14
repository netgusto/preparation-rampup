package main

import (
	"fmt"
	"io/ioutil"
	"sync"
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

	pingStream := make(chan Ping)
	urlStream := make(chan string)
	done := make(chan bool)

	wg := &sync.WaitGroup{}

	for w := 0; w < 4; w++ {
		go worker(w, urlStream, pingStream, wg, nbTriesPerURL, urlGetter)
	}

	go func() {
		for ping := range pingStream {
			fmt.Printf("%+v\n", ping)
			wg.Done()
		}

		done <- true
	}()

	for round := 0; round < 10; round++ {
		for _, url := range urls {
			urlStream <- url
		}
	}

	wg.Wait()
	close(pingStream)
	<-done

	fmt.Println("DONE!")
}

func worker(id int, urls <-chan string, pings chan<- Ping, wg *sync.WaitGroup, nbTriesPerURL int, urlGetter URLGetter) {
	for url := range urls {
		wg.Add(1)
		pings <- probeURL(url, nbTriesPerURL, urlGetter)
	}
}
