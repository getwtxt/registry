package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// AddUser inserts a new user into the calling Index registry object.
func (index *Index) AddUser(nick, urls string, ipaddr net.IP, statuses TimeMap) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("can't add user to uninitialized index")

	} else if nick == "" || urls == "" {
		// Check that the request is valid.
		return fmt.Errorf("both URL and Nick must be specified")

	} else if !strings.HasPrefix(urls, "http") {
		return fmt.Errorf("invalid URL: %v", urls)
	}

	// Check if a user already exists in the index.
	index.Mu.RLock()
	if _, ok := index.Reg[urls]; ok {
		index.Mu.RUnlock()
		return fmt.Errorf("user %v already exists", urls)
	}
	index.Mu.RUnlock()

	// Get the time as both a standard time.Time and as
	// an RFC3339-formatted timestamp. This will be used
	// for Data.Date and Data.APIdate to record when the
	// user was added to the index. Ignoring the error
	// because of the near-nil possibility of it happening
	thetime := time.Now()

	// Acquire a write lock and load the user data into
	// our index.
	index.Mu.Lock()
	index.Reg[urls] = &Data{
		Mu:     sync.RWMutex{},
		Nick:   nick,
		URL:    urls,
		IP:     ipaddr,
		Date:   thetime.Format(time.RFC3339),
		Status: statuses}
	index.Mu.Unlock()

	return nil
}

// DelUser removes a user and all associated data from
// the calling Index registry object.
func (index *Index) DelUser(urls string) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("can't delete user from empty index")

	} else if urls == "" {
		// Check that the request is valid.
		return fmt.Errorf("can't delete blank user")

	} else if !strings.HasPrefix(urls, "http") {
		// Check that we were provided a URL
		return fmt.Errorf("invalid URL: %v", urls)
	}

	// Check if the user exists in the index. If
	// they don't, we can't remove them.
	index.Mu.RLock()
	if _, ok := index.Reg[urls]; !ok {
		index.Mu.RUnlock()
		return fmt.Errorf("can't delete user %v, user doesn't exist", urls)
	}
	index.Mu.RUnlock()

	// Acquire a write lock and delete the user from
	// the index.
	index.Mu.Lock()
	delete(index.Reg, urls)
	index.Mu.Unlock()

	return nil
}

// UpdateUser adds new statuses to the user's entry in the registry.
func (index *Index) UpdateUser(urls string) error {
	// fetch the twtxt file data
	out, registry, err := GetTwtxt(urls)
	if err != nil {
		return err
	}

	// if we've somehow fetched a registry's data, error out
	if registry {
		return fmt.Errorf("attempting to update registry URL - users should be updated individually")
	}

	index.Mu.RLock()
	user := index.Reg[urls]
	index.Mu.RUnlock()
	nick := user.Nick

	// update the user's entry in the Index
	data, err := ParseUserTwtxt(out, nick, urls)
	if err != nil {
		return err
	}
	index.Mu.Lock()
	tmp := index.Reg[urls]
	for i, e := range data {
		tmp.Status[i] = e
	}
	index.Mu.Unlock()

	return nil
}

// ScrapeRemoteRegistry adds the users who are available via remote registry
func (index *Index) ScrapeRemoteRegistry(urls string) error {
	//fetch the remote registry's entries
	out, registry, err := GetTwtxt(urls)
	if err != nil {
		return err
	}

	// if we're working with an individual's twtxt file, error out
	if !registry {
		return fmt.Errorf("can't add single user via call to AddRemoteRegistry")
	}

	// parse the registry's data and add to our Index
	data, err := ParseRegistryTwtxt(out)
	if err != nil {
		return err
	}

	// only add new users so we don't overwrite data
	// we already have (and lose statuses, etc)
	index.Mu.Lock()
	for _, e := range data {
		if _, ok := index.Reg[e.URL]; !ok {
			index.Reg[e.URL] = e
		}
	}
	index.Mu.Unlock()

	return nil
}

// GetUserStatuses returns a TimeMap containing single user's statuses
func (index *Index) GetUserStatuses(urls string) (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil {
		return nil, fmt.Errorf("can't get statuses from an empty index")

	} else if urls == "" {
		// Check that the request is valid.
		return nil, fmt.Errorf("can't retrieve statuses of blank user")

	} else if !strings.HasPrefix(urls, "http") {
		// Check that we were provided a URL
		return nil, fmt.Errorf("invalid URL: %v", urls)
	}

	// Check if the user is in the index. If they
	// aren't, we can't return their statuses.
	index.Mu.RLock()
	if _, ok := index.Reg[urls]; !ok {
		index.Mu.RUnlock()
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}

	// Pull the user's statuses from the index.
	status := index.Reg[urls].Status
	index.Mu.RUnlock()

	return status, nil
}

// GetStatuses returns a TimeMap containing all statuses
// from all users in the calling Index registry object.
func (index *Index) GetStatuses() (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil {
		return nil, fmt.Errorf("can't get statuses from an empty index")
	}

	// Initialize a new TimeMap in which we'll
	// store the statuses.
	statuses := NewTimeMap()

	// For each user, assign each status to
	// our aggregate TimeMap.
	index.Mu.RLock()
	for _, v := range index.Reg {
		if v.Status == nil || len(v.Status) == 0 {
			// Skip a user's statuses if the map is uninitialized or zero length
			continue
		}
		for a, b := range v.Status {
			if _, ok := v.Status[a]; ok {
				statuses[a] = b
			}
		}
	}
	index.Mu.RUnlock()

	return statuses, nil
}
