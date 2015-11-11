package scrape

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

type empty struct{}

var sem chan empty

// Initialize semaphore
func init() {
	sem = make(chan empty, 20)
}

// Process requests to course API. Use semaphore to limit the number of
// concurrent connections to the API.
func GetXML(url string) ([]byte, error) {

	// Acquire
	e := empty{}
	sem <- e

	defer func() {
		<-sem
	}()

	return getXML(url)
}

// Make request to url and return XML at that url.
func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error("HTTP GET Failed:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Warn("oops! recieved: ", resp.StatusCode)
		time.Sleep(2 * time.Second)
		resp.Body.Close()
		return getXML(url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body: ", err)
		return nil, err
	}

	return body, nil
}
