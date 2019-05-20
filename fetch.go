package registry // import "github.com/getwtxt/registry"

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// GetTwtxt fetches the raw twtxt file data from the user's
// provided URL, after validating the URL.
func GetTwtxt(urls string) ([]byte, bool, error) {

	// Check that we were provided a valid
	// URL in the first place.
	if !strings.HasPrefix(urls, "http") {
		return nil, fmt.Errorf("invalid twtxt file url: %v", urls)
	}

	// Request the data
	req, err := http.Get(urls)
	if err != nil {
		return nil, fmt.Errorf("couldn't get %v: %v", urls, err)
	}
	defer func() {
		err := req.Body.Close()
		if err != nil {
			log.Printf("Couldn't close response body for %v: %v\n", urls, err)
		}
	}()

	// Verify that we've received text-only content
	// and not something else.
	var textplain bool
	for _, v := range req.Header["Content-Type"] {
		if strings.Contains(v, "text/plain") {
			textplain = true
			break
		}
	}
	if !textplain {
		return nil, false, fmt.Errorf("received non-text/plain response body from %v", urls)
	}

	// Make sure the request returned a 200
	if req.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("didn't get 200 from remote server, received %v: %v", req.StatusCode, urls)
	}

	// Pull the response body into a variable
	twtxt, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, false, fmt.Errorf("error reading response body from %v: %v", urls, err)
	}

	// Signal that we're adding another twtxt registry as a "user"
	if strings.HasSuffix(urls, "/tweets") {
		return twtxt, true, nil
	}

	return twtxt, false, nil
}

// ParseTwtxt takes a fetched twtxt file in the form of
// a slice of bytes, parses it, and returns it as a
// TimeMap. The output may then be passed to AddUser()
func ParseTwtxt(twtxt []byte, otherRegistryOutput bool) (TimeMap, []error) {
	// Store timestamp parsing errors in a slice
	// of errors.
	var erz []error

	// Make sure we actually have something to parse
	if len(twtxt) == 0 {
		return nil, append(erz, fmt.Errorf("received no data"))
	}

	// Set everything up to parse the twtxt file
	reader := bytes.NewReader(twtxt)
	scanner := bufio.NewScanner(reader)
	timemap := NewTimeMap()

	// Scan the data by linebreak
	for scanner.Scan() {
		thetime := time.Time{}
		nopadding := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(nopadding, "#") || nopadding == "" {
			continue
		}

		// Split the twtxt file into columns by tabs
		columns := strings.Split(nopadding, "\t")
		if len(columns) != 2 && !otherRegistryOutput {
			return nil, append(erz, fmt.Errorf("improperly formatted data"))
		}
		if len(columns) != 4 && otherRegistryOutput {
			return nil, append(erz, fmt.Errorf("improperly formatted data"))
		}

		// Take the RFC3339 date in the third column
		// and convert it into a standard time.Time.
		// If there was a parsing error, keep going,
		// but take note.
		var err error
		if !otherRegistryOutput {
			err = thetime.UnmarshalText([]byte(columns[0]))
		} else {
			err = thetime.UnmarshalText([]byte(columns[2]))
		}
		if err != nil {
			erz = append(erz, fmt.Errorf("unable to retrieve date: %v", err))
		}

		// Add the status to the TimeMap
		timemap[thetime] = scanner.Text()
	}
	return timemap, erz
}
