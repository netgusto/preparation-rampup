package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
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
	run(urlGetter, urls)

	fmt.Println("DONE!")
}

func run(urlGetter URLGetter, urls []string) {
	pingStream := make(chan Ping)

	go func() {
		// consume ping stream
		for ping := range pingStream {
			fmt.Printf("%+v\n", ping)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.TODO())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancelFunc()
		println("CANCELLED!")
	}()

	probeLoop(ctx, urlGetter, urls, 4, 10, pingStream)
}
