package main

import (
	"context"
)

func probeLoop(ctx context.Context, urlGetter URLGetter, urls []string, nbWorkers int, nbRounds int, pingStream chan<- Ping) {

	workersCtx, cancelWorkers := context.WithCancel(context.Background())
	urlStream := make(chan string)

	for w := 0; w < nbWorkers; w++ {
		go worker(workersCtx, w, urlStream, pingStream, nbTriesPerURL, urlGetter)
	}

	for {
		for _, url := range urls {
			select {
			case <-ctx.Done():
				cancelWorkers()
				return
			default:
				urlStream <- url
			}
		}
	}
}

func worker(ctx context.Context, id int, urlStream <-chan string, pingStream chan<- Ping, nbTriesPerURL int, urlGetter URLGetter) {
	for url := range urlStream {
		ping := probeURL(url, nbTriesPerURL, urlGetter)

		select {
		case <-ctx.Done():
			return
		default:
			pingStream <- ping
		}
	}
}
