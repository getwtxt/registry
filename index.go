// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"log"
	"time"
)

// AddUser inserts a new user into the index. The *Data struct
// contains the nickname and the time the user was added.
// TODO: Tie this to GetTwtxt / ParseTwtxt
func (index UserIndex) AddUser(nick string, url string) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("index hasn't been initialized")
	}

	// Check that the request is valid.
	if nick == "" || url == "" {
		return fmt.Errorf("both URL and Nick must be specified")
	}

	// Use the double-return-value property of map lookups
	// to check if a user already exists in the index.
	imutex.RLock()
	if _, ok := index[url]; ok {
		imutex.RUnlock()
		log.Printf("User %v can't be added - already exists.\n", url)
		return fmt.Errorf("user %v already exists", url)
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

	// Acquire a write lock and load the user data into
	// our index.
	imutex.Lock()
	index[url] = &Data{Nick: nick, Date: thetime, APIdate: rfc3339date}
	imutex.Unlock()

	return nil
}

// DelUser removes a user from the index completely.
func (index UserIndex) DelUser(url string) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("index hasn't been initialized")
	}

	// Check that the request is valid.
	if url == "" {
		return fmt.Errorf("can't delete blank user")
	}

	// Use the double-return-value property of maps
	// to check if the user exists in the index. If
	// they don't, we can't remove them.
	imutex.RLock()
	if _, ok := index[url]; !ok {
		imutex.RUnlock()
		return fmt.Errorf("can't delete user %v, user doesn't exist", url)
	}
	imutex.RUnlock()

	// Acquire a write lock and delete the user from
	// the index.
	imutex.Lock()
	delete(index, url)
	imutex.Unlock()

	return nil
}

// GetUserStatuses returns a TimeMap containing a user's statuses
func (index UserIndex) GetUserStatuses(url string) (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil {
		return nil, fmt.Errorf("index hasn't been initialized")
	}

	// Check that the request is valid.
	if url == "" {
		return nil, fmt.Errorf("can't retrieve statuses of blank user")
	}

	// Use the double-return-value property of maps
	// to check if the user is in the index. If they
	// aren't, we can't return their statuses.
	imutex.RLock()
	if _, ok := index[url]; !ok {
		imutex.RUnlock()
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}

	// Pull the user's statuses from the index.
	status := index[url].Status
	imutex.RUnlock()

	return status, nil
}

// GetStatuses returns a TimeMap containing all statuses
// in the index.
func (index UserIndex) GetStatuses() (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil {
		return nil, fmt.Errorf("can't retrieve statuses from empty index")
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
