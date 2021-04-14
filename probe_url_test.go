package main

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

type URLGetterMockResponse struct {
	StatusCode *int
	Duration   time.Duration
	Error      error
}

type URLGetterMock struct {
	Responses  []URLGetterMockResponse
	CurrentTry int
}

func (ugmock *URLGetterMock) GetURL(url string) (*int, time.Duration, error) {
	try := ugmock.CurrentTry
	ugmock.CurrentTry++

	return ugmock.Responses[try].StatusCode, ugmock.Responses[try].Duration, ugmock.Responses[try].Error
}

var statusOK = http.StatusOK

func Test_ProbeURL_with_working_url(t *testing.T) {

	maxTries := 3

	mockURLGetter := &URLGetterMock{
		Responses: []URLGetterMockResponse{
			{StatusCode: &statusOK, Duration: time.Millisecond * 150, Error: nil},
		},
	}

	actualPing := probeURL("https://whatever", maxTries, mockURLGetter)

	if actualPing.StatusCode == nil || *actualPing.StatusCode != http.StatusOK {
		t.Fatalf(`Expected the test to yield HTTP 200 but was %v`, actualPing.StatusCode)
	}

	if actualPing.TransportError != nil {
		t.Fatalf(`Expected the probing to not error but got %v`, actualPing.TransportError)
	}

	if actualPing.Tries != 1 {
		t.Fatalf(`Expected the probing to be done in 1 try, but was tried %d times`, actualPing.Tries)
	}

	if actualPing.URL != "https://whatever" {
		t.Fatalf(`Got invalid URL %s`, actualPing.URL)
	}

	if actualPing.Duration != time.Millisecond*150 {
		t.Fatalf(`Invalid duration %v`, actualPing.Duration)
	}

	if time.Since(actualPing.Moment) > time.Second*10 {
		t.Fatalf(`Moment looks broken %v`, actualPing.Moment)
	}
}

func Test_ProbeURL_with_unreachable_url(t *testing.T) {

	maxTries := 3
	urlGetErr := errors.New("TEST ERROR: Could not get the URL!")

	mockURLGetter := &URLGetterMock{
		Responses: []URLGetterMockResponse{
			{StatusCode: nil, Duration: time.Millisecond * 150, Error: urlGetErr},
			{StatusCode: nil, Duration: time.Millisecond * 200, Error: urlGetErr},
			{StatusCode: nil, Duration: time.Millisecond * 250, Error: urlGetErr},
		},
	}

	actualPing := probeURL("https://whatever", maxTries, mockURLGetter)

	if actualPing.StatusCode != nil {
		t.Fatalf(`Expected the test to yield nil HTTP status code but was %v`, actualPing.StatusCode)
	}

	if actualPing.TransportError != urlGetErr {
		t.Fatalf(`Invalid error`)
	}

	if actualPing.Tries != maxTries {
		t.Fatalf(`Expected the probing to be tried %d times, but was tried %d times`, maxTries, actualPing.Tries)
	}

	if actualPing.URL != "https://whatever" {
		t.Errorf(`Got invalid URL %s`, actualPing.URL)
	}
}

func Test_ProbeURL_ok_on_last_try(t *testing.T) {

	maxTries := 3
	urlGetErr := errors.New("TEST ERROR: Could not get the URL!")

	mockURLGetter := &URLGetterMock{
		Responses: []URLGetterMockResponse{
			{StatusCode: nil, Duration: time.Millisecond * 150, Error: urlGetErr},
			{StatusCode: nil, Duration: time.Millisecond * 200, Error: urlGetErr},
			{StatusCode: &statusOK, Duration: time.Millisecond * 250, Error: nil},
		},
	}

	actualPing := probeURL("https://whatever", maxTries, mockURLGetter)

	if actualPing.StatusCode == nil || *actualPing.StatusCode != http.StatusOK {
		t.Fatalf(`Expected the test to yield HTTP 200 but was %v`, actualPing.StatusCode)
	}

	if actualPing.TransportError != nil {
		t.Fatalf(`Unexpected error`)
	}

	if actualPing.Tries != maxTries {
		t.Fatalf(`Expected the probing to be tried %d times, but was tried %d times`, maxTries, actualPing.Tries)
	}

	if actualPing.URL != "https://whatever" {
		t.Errorf(`Got invalid URL %s`, actualPing.URL)
	}
}
