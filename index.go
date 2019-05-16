// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// AddUser inserts a new user into the index. The *Data struct
// contains the nickname and the time the user was added.
// TODO: Tie this to GetTwtxt / ParseTwtxt
func (index UserIndex) AddUser(nick string, urls string) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("can't add user to uninitialized index")
	}

	// Check that the request is valid.
	if nick == "" || urls == "" {
		return fmt.Errorf("both URL and Nick must be specified")
	} else if !strings.HasPrefix(urls, "http") {
		return fmt.Errorf("invalid URL: %v", urls)
	}

	// Use the double-return-value property of map lookups
	// to check if a user already exists in the index.
	imutex.RLock()
	if _, ok := index[urls]; ok {
		imutex.RUnlock()
		log.Printf("User %v can't be added - already exists.\n", urls)
		return fmt.Errorf("user %v already exists", urls)
	}
	imutex.RUnlock()

	// Get the time as both a standard time.Time and as
	// an RFC3339-formatted timestamp. This will be used
	// for Data.Date and Data.APIdate to record when the
	// user was added to the index.
	thetime := time.Now()
	rfc3339date, err := thetime.MarshalText()
	if err != nil {
		log.Printf("Error formatting user add time as RFC3339: %v\n", err)
	}

	// Retrieve and parse the user's twtxt status file
	rawdata, err := GetTwtxt(urls)
	if err != nil {
		return fmt.Errorf("error retrieving twtxt status file: %v", err)
	}
	parsed, errz := ParseTwtxt(rawdata)
	if err != nil {
		out := "error(s) parsing twtxt status file: "
		for _, e := range errz {
			out += fmt.Sprintf("%v, ", e)
		}
		return fmt.Errorf("%v", out)
	}

	// Acquire a write lock and load the user data into
	// our index.
	imutex.Lock()
	index[urls] = &Data{Nick: nick, Date: thetime, APIdate: rfc3339date, Status: parsed}
	imutex.Unlock()

	return nil
}

// DelUser removes a user from the index completely.
func (index UserIndex) DelUser(urls string) error {

	// Check that we have an initialized index.
	if index == nil || len(index) == 0 {
		return fmt.Errorf("can't delete user from empty index")
	}

	// Check that the request is valid.
	if urls == "" {
		return fmt.Errorf("can't delete blank user")
	} else if !strings.HasPrefix(urls, "http") {
		return fmt.Errorf("invalid URL: %v", urls)
	}

	// Use the double-return-value property of maps
	// to check if the user exists in the index. If
	// they don't, we can't remove them.
	imutex.RLock()
	if _, ok := index[urls]; !ok {
		imutex.RUnlock()
		return fmt.Errorf("can't delete user %v, user doesn't exist", urls)
	}
	imutex.RUnlock()

	// Acquire a write lock and delete the user from
	// the index.
	imutex.Lock()
	delete(index, urls)
	imutex.Unlock()

	return nil
}

// GetUserStatuses returns a TimeMap containing a user's statuses
func (index UserIndex) GetUserStatuses(urls string) (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil || len(index) == 0 {
		return nil, fmt.Errorf("can't get statuses from an empty index")
	}

	// Check that the request is valid.
	if urls == "" {
		return nil, fmt.Errorf("can't retrieve statuses of blank user")
	} else if !strings.HasPrefix(urls, "http") {
		return nil, fmt.Errorf("invalid URL: %v", urls)
	}

	// Use the double-return-value property of maps
	// to check if the user is in the index. If they
	// aren't, we can't return their statuses.
	imutex.RLock()
	if _, ok := index[urls]; !ok {
		imutex.RUnlock()
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}

	// Pull the user's statuses from the index.
	status := index[urls].Status
	imutex.RUnlock()

	return status, nil
}

// GetStatuses returns a TimeMap containing all statuses
// in the index.
func (index UserIndex) GetStatuses() (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil || len(index) == 0 {
		return nil, fmt.Errorf("can't get statuses from an empty index")
	}

	// Initialize a new TimeMap in which we'll
	// store the statuses.
	statuses := NewTimeMap()

	// For each user, assign each status to
	// our aggregate TimeMap. This needs to
	// be refactored badly. as it's O(n^2)
	// and probably doesn't need to be.
	imutex.RLock()
	for _, v := range index {
		for a, b := range v.Status {
			statuses[a] = b
		}
	}
	imutex.RUnlock()

	return statuses, nil
}
