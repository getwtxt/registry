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
// GetTwtxt should be passed to either ParseUserTwtxt or
// ParseRegistryTwtxt, respectively.
func GetTwtxt(urlKey string) ([]byte, bool, error) {

	// Check that we were provided a valid
	// URL in the first place.
	if !strings.HasPrefix(urlKey, "http") {
		return nil, false, fmt.Errorf("invalid twtxt file url: %v", urlKey)
	}

	// Set the timeout for all requests
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Craft a request
	var b []byte
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest("GET", urlKey, buf)
	if err != nil {
		return nil, false, err
	}

	// Request the data
	res, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("couldn't get %v: %v", urlKey, err)
	}

	defer res.Body.Close()

	// Verify that we've received text-only content
	// and not something else.
	var textPlain bool
	for _, v := range res.Header["Content-Type"] {
		if strings.Contains(v, "text/plain") {
			textPlain = true
			break
		}
	}
	if !textPlain {
		return nil, false, fmt.Errorf("received non-text/plain response body from %v", urlKey)
	}

	// Make sure the request returned a 200
	if res.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("didn't get 200 from remote server, received %v: %v", res.StatusCode, urlKey)
	}

	// Pull the response body into a variable
	twtxt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, false, fmt.Errorf("error reading response body from %v: %v", urlKey, err)
	}

	// Signal that we're adding another twtxt registry as a "user"
	if strings.HasSuffix(urlKey, "/api/plain/tweets") {
		return twtxt, true, nil
	}

	return twtxt, false, nil
}

// ParseUserTwtxt takes a fetched twtxt file in the form of
// a slice of bytes, parses it, and returns it as a
// TimeMap. The output may then be passed to Index.AddUser()
func ParseUserTwtxt(twtxt []byte, nickname, urlKey string) (TimeMap, error) {
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
		thetime, err := time.Parse(time.RFC3339, columns[0])
		if err != nil {
			erz = append(erz, []byte(fmt.Sprintf("unable to retrieve date: %v\n", err))...)
		}

		// Add the status to the TimeMap
		timemap[thetime] = nickname + "\t" + urlKey + "\t" + nopadding
	}

	if len(erz) == 0 {
		return timemap, nil
	}
	return timemap, fmt.Errorf("%v", erz)
}

// ParseRegistryTwtxt takes output from a remote registry and outputs
// the accessible user data via a slice of Users.
func ParseRegistryTwtxt(twtxt []byte) ([]*User, error) {

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
	userdata := []*User{}

	// Scan the data by linebreak
	for scanner.Scan() {

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
		thetime, err := time.Parse(time.RFC3339, columns[2])
		if err != nil {
			erz = append(erz, []byte(fmt.Sprintf("%v\n", err))...)
			continue
		}

		parsednickname := columns[0]
		dataIndex := 0
		inIndex := false
		parsedurl := columns[1]

		for i, e := range userdata {
			if e.Nick == parsednickname || e.URL == parsedurl {
				dataIndex = i
				inIndex = true
				break
			}
		}

		if inIndex {
			tmp := userdata[dataIndex]
			tmp.Status[thetime] = nopadding
			userdata[dataIndex] = tmp
		} else {
			// If the user hasn't been seen before,
			// create a new Data object
			timeNowRFC := time.Now().Format(time.RFC3339)
			if err != nil {
				erz = append(erz, []byte(fmt.Sprintf("%v\n", err))...)
			}

			tmp := &User{
				Mu:   sync.RWMutex{},
				Nick: parsednickname,
				URL:  parsedurl,
				Date: timeNowRFC,
				Status: TimeMap{
					thetime: nopadding,
				},
			}

			userdata = append(userdata, tmp)
		}

		inIndex = false
	}

	return userdata, fmt.Errorf("%v", erz)
}
