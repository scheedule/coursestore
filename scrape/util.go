package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type empty struct{}

var sem chan empty

func init() {
	sem = make(chan empty, 20)
}

func GetXML(url string) ([]byte, error) {

	// Acquire
	e := empty{}
	sem <- e

	defer func() {
		<-sem
	}()

	return getXML(url)
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Oops! recieved: ", resp.StatusCode)
		time.Sleep(2 * time.Second)
		resp.Body.Close()
		return getXML(url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
