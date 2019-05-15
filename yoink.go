package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http" // GetTwtxt fetches the raw twtxt.txt file from the remote server
	"net/url"
)

// GetTwtxt fetches the raw twtxt.txt data from the user's
// provided URL, after validating the URL.
func GetTwtxt(urls string) ([]byte, error) {
	_, err := url.Parse(urls)
	if err != nil {
		return nil, fmt.Errorf("invalid twtxt.txt url: %v, %v", urls, err)
	}

	req, err := http.Get(urls)
	if err == nil {
		defer func() {
			err := req.Body.Close()
			if err != nil {
				log.Printf("Couldn't close response body for %v: %v\n", urls, err)
			}
		}()
	} else if err != nil {
		return nil, fmt.Errorf("couldn't get %v: %v", urls, err)
	} else if req.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("didn't get 200 from remote server, received %v: %v", req.StatusCode, urls)
	}

	twtxt, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body from %v: %v", urls, err)
	}

	return twtxt, nil
}
