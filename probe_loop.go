package main

import "sync"

func probeLoop(urlGetter URLGetter, urls []string, nbWorkers int, nbRounds int, pingStream chan<- Ping) {

	wg := &sync.WaitGroup{}

	urlStream := make(chan string)

	for w := 0; w < nbWorkers; w++ {
		go worker(w, wg, urlStream, pingStream, nbTriesPerURL, urlGetter)
	}

	for round := 0; round <= nbRounds; round++ {
		for _, url := range urls {
			wg.Add(1)
			urlStream <- url
		}
	}

	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup, urlStream <-chan string, pingStream chan<- Ping, nbTriesPerURL int, urlGetter URLGetter) {
	for url := range urlStream {
		pingStream <- probeURL(url, nbTriesPerURL, urlGetter)
		wg.Done()
	}
}
