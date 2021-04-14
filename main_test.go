package main

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_ProbeURL_with_working_url(t *testing.T) {

	maxTries := 3
	start := time.Now()

	actualPing := probeURL("https://www.algolia.com", maxTries)

	if actualPing.StatusCode == nil || *actualPing.StatusCode != http.StatusOK {
		t.Fatalf(`Expected the test to yield HTTP 200 but was %v`, actualPing.StatusCode)
	}

	if actualPing.TransportError != nil {
		t.Fatalf(`Expected the probing to not error but got %v`, actualPing.TransportError)
	}

	if actualPing.Tries != 1 {
		t.Fatalf(`Expected the probing to be done in 1 try, but was tried %d times`, actualPing.Tries)
	}

	if actualPing.URL != "https://www.algolia.com" {
		t.Fatalf(`Got invalid URL %s`, actualPing.URL)
	}

	if actualPing.Duration > time.Since(start) || actualPing.Duration == 0 {
		t.Fatalf(`Duration looks broken %v`, actualPing.Duration)
	}

	if time.Since(actualPing.Moment) > time.Second*10 {
		t.Fatalf(`Moment looks broken %v`, actualPing.Moment)
	}
}

func Test_ProbeURL_with_unresolvable_url(t *testing.T) {

	maxTries := 3
	actualPing := probeURL("https://this-domain-does-not-exist", maxTries)

	if actualPing.StatusCode != nil {
		t.Fatalf(`Expected the test to yield nil HTTP status code 200 but was %v`, actualPing.StatusCode)
	}

	if actualPing.TransportError == nil {
		t.Fatalf(`Expected the probing to error but got nil`)
	}

	if !strings.Contains(actualPing.TransportError.Error(), "no such host") {
		t.Fatalf(`Invalid error, got %s`, actualPing.TransportError.Error())
	}

	if actualPing.Tries != maxTries {
		t.Fatalf(`Expected the probing to be tried %d times, but was tried %d times`, maxTries, actualPing.Tries)
	}

	if actualPing.URL != "https://this-domain-does-not-exist" {
		t.Errorf(`Got invalid URL %s`, actualPing.URL)
	}
}
