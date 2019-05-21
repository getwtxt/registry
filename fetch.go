package registry // import "github.com/getwtxt/registry"

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// GetTwtxt fetches the raw twtxt file data from the user's
// provided URL, after validating the URL. If the returned
// boolean value is false, the fetched URL is a single user's
// twtxt file. If true, the fetched URL is the output of
// another registry's /api/plain/tweets. The output of
// GetTwtxt should be passed to either ParseTwtxt or
// ParseRegistryTwtxt, respectively.
func GetTwtxt(urls string) ([]byte, bool, error) {

	// Check that we were provided a valid
	// URL in the first place.
	if !strings.HasPrefix(urls, "http") {
		return nil, false, fmt.Errorf("invalid twtxt file url: %v", urls)
	}

	// Request the data
	req, err := http.Get(urls)
	if err != nil {
		return nil, false, fmt.Errorf("couldn't get %v: %v", urls, err)
	}

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
	if strings.HasSuffix(urls, "/api/plain/tweets") {
		return twtxt, true, nil
	}

	return twtxt, false, nil
}

// ParseUserTwtxt takes a fetched twtxt file in the form of
// a slice of bytes, parses it, and returns it as a
// TimeMap. The output may then be passed to AddUser()
func ParseUserTwtxt(twtxt []byte) (TimeMap, error) {
	// Store timestamp parsing errors in a slice
	// of errors.
	var erz []byte

	// Make sure we actually have something to parse
	if len(twtxt) == 0 {
		return nil, fmt.Errorf("no data to parse in twtxt file")
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
		if len(columns) != 2 {
			return nil, fmt.Errorf("improperly formatted data in twtxt file")
		}

		// Take the RFC3339 date in the third column
		// and convert it into a standard time.Time.
		// If there was a parsing error, keep going,
		// but take note.
		err := thetime.UnmarshalText([]byte(columns[0]))
		if err != nil {
			erz = append(erz, []byte(fmt.Sprintf("unable to retrieve date: %v\n", err))...)
		}

		// Add the status to the TimeMap
		timemap[thetime] = scanner.Text()
	}
	if len(erz) == 0 {
		return timemap, nil
	}
	return timemap, fmt.Errorf("%v", erz)
}

// ParseRegistryTwtxt takes output from another registry and outputs it
// via a slice of Data objects.
func ParseRegistryTwtxt(twtxt []byte) ([]*Data, error) {

	// Store timestamp parsing errors in a slice
	// of errors.
	var erz []byte

	// Make sure we actually have something to parse
	if len(twtxt) == 0 {
		return nil, fmt.Errorf("received no data")
	}

	// Set everything up to parse the twtxt file
	reader := bytes.NewReader(twtxt)
	scanner := bufio.NewScanner(reader)
	userdata := []*Data{}

	// Scan the data by linebreak
	for scanner.Scan() {

		thetime := time.Time{}
		nopadding := strings.TrimSpace(scanner.Text())

		// check if we've happened upon a comment or a blank line
		if strings.HasPrefix(nopadding, "#") || nopadding == "" {
			continue
		}

		// Split the twtxt file into columns by tabs
		columns := strings.Split(nopadding, "\t")
		if len(columns) != 4 {
			return nil, fmt.Errorf("improperly formatted data")
		}

		// Take the RFC3339 date in the third column
		// and convert it into a standard time.Time.
		// If there was a parsing error, keep going
		// and skip that status.
		err := thetime.UnmarshalText([]byte(columns[2]))
		if err != nil {
			erz = append(erz, []byte(fmt.Sprintf("%v\n", err))...)
			continue
		}

		parsednick := columns[0]
		dataIndex := 0
		inIndex := false
		parsedurl := columns[1]

		for i, e := range userdata {
			if e.Nick == parsednick || e.URL == parsedurl {
				dataIndex = i
				inIndex = true
				break
			}
		}

		if inIndex {
			tmp := userdata[dataIndex]
			tmp.Status[thetime] = columns[2] + "\t" + columns[3]
			userdata[dataIndex] = tmp
		} else {
			// If the user hasn't been seen before,
			// create a new Data object
			timeNow := time.Now()

			timeNowRFC, err := timeNow.MarshalText()
			if err != nil {
				erz = append(erz, []byte(fmt.Sprintf("%v\n", err))...)
			}

			tmp := &Data{
				Mu:      sync.RWMutex{},
				Nick:    parsednick,
				URL:     parsedurl,
				Date:    timeNow,
				APIdate: timeNowRFC,
				Status: TimeMap{
					thetime: columns[2] + "\t" + columns[3],
				},
			}

			userdata = append(userdata, tmp)
		}

		inIndex = false
	}

	return userdata, fmt.Errorf("%v", erz)
}
