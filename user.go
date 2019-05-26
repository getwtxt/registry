package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// AddUser inserts a new user into the Index.
func (index *Index) AddUser(nickname, urlKey string, ipAddress net.IP, statuses TimeMap) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("can't add user to uninitialized index")

	} else if nickname == "" || urlKey == "" {
		// Check that the request is valid.
		return fmt.Errorf("both URL and Nick must be specified")

	} else if !strings.HasPrefix(urlKey, "http") {
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

	// Check if a user already exists in the index.
	index.Mu.RLock()
	if _, ok := index.Reg[urlKey]; ok {
		index.Mu.RUnlock()
		return fmt.Errorf("user %v already exists", urlKey)
	}
	index.Mu.RUnlock()

	// Get the time as both a standard time.Time and as
	// an RFC3339-formatted timestamp. This will be used
	// for User.Date and User.APIdate to record when the
	// user was added to the index. Ignoring the error
	// because of the near-nil possibility of it happening
	thetime := time.Now()

	// Acquire a write lock and load the user data into
	// our index.
	index.Mu.Lock()
	index.Reg[urlKey] = &User{
		Mu:     sync.RWMutex{},
		Nick:   nickname,
		URL:    urlKey,
		IP:     ipAddress,
		Date:   thetime.Format(time.RFC3339),
		Status: statuses}
	index.Mu.Unlock()

	return nil
}

// Push inserts a given User into an Index. The User
// being pushed need only have the URL field filled.
// All other fields may be empty.
// This can be destructive: an existing User in the
// Index will be overwritten if its User.URL is the
// same as the User.URL being pushed.
func (index *Index) Push(user *User) error {
	if user == nil {
		return fmt.Errorf("can't push nil data to index")
	}
	if index == nil || index.Reg == nil {
		return fmt.Errorf("can't push data to index: index uninitialized")
	}
	user.Mu.RLock()
	if user.URL == "" {
		user.Mu.RUnlock()
		return fmt.Errorf("can't push data to index: missing URL for key")
	}
	urlKey := user.URL
	index.Mu.Lock()
	index.Reg[urlKey] = user
	index.Mu.Unlock()
	user.Mu.RUnlock()

	return nil
}

// Pop pulls the User associated with the
// provided URL key from the Index. It then
// removes the User from the Index and returns
// the User.
func (index *Index) Pop(urlKey string) (*User, error) {
	if index == nil {
		return nil, fmt.Errorf("can't pop from nil index")
	}
	if urlKey == "" {
		return nil, fmt.Errorf("can't pop unless provided a key")
	}

	index.Mu.RLock()
	if _, ok := index.Reg[urlKey]; !ok {
		index.Mu.RUnlock()
		return nil, fmt.Errorf("provided url key doesn't exist in index")
	}

	index.Reg[urlKey].Mu.RLock()
	userUser := index.Reg[urlKey]
	index.Reg[urlKey].Mu.RUnlock()
	index.Mu.RUnlock()

	index.Mu.Lock()
	index.Reg[urlKey].Mu.Lock()
	delete(index.Reg, urlKey)
	index.Mu.Unlock()

	return userUser, nil
}

// DelUser removes a user and all associated data from
// the Index.
func (index *Index) DelUser(urlKey string) error {

	// Check that we have an initialized index.
	if index == nil {
		return fmt.Errorf("can't delete user from empty index")

	} else if urlKey == "" {
		// Check that the request is valid.
		return fmt.Errorf("can't delete blank user")

	} else if !strings.HasPrefix(urlKey, "http") {
		// Check that we were provided a URL
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

	// Check if the user exists in the index. If
	// they don't, we can't remove them.
	index.Mu.RLock()
	if _, ok := index.Reg[urlKey]; !ok {
		index.Mu.RUnlock()
		return fmt.Errorf("can't delete user %v, user doesn't exist", urlKey)
	}
	index.Mu.RUnlock()

	// Acquire a write lock and delete the user from
	// the index. The User mutex is never unlocked because
	// the User is deleted. It is only acquired to
	// prevent a panic if another thread is reading/writing
	// to the user.
	index.Mu.Lock()
	index.Reg[urlKey].Mu.Lock()
	delete(index.Reg, urlKey)
	index.Mu.Unlock()

	return nil
}

// UpdateUser scrapes an existing user's remote twtxt.txt
// file. Any new statuses are added to the user's entry
// in the Index.
func (index *Index) UpdateUser(urlKey string) error {
	// fetch the twtxt file data
	out, registry, err := GetTwtxt(urlKey)
	if err != nil {
		return err
	}

	// if we've somehow fetched a registry's data, error out
	if registry {
		return fmt.Errorf("attempting to update registry URL - users should be updated individually")
	}

	index.Mu.RLock()
	user := index.Reg[urlKey]
	index.Mu.RUnlock()
	user.Mu.RLock()
	nick := user.Nick
	user.Mu.RUnlock()

	// update the user's entry in the Index
	data, err := ParseUserTwtxt(out, nick, urlKey)
	if err != nil {
		return err
	}
	index.Mu.Lock()
	tmp := index.Reg[urlKey]
	tmp.Mu.Lock()
	for i, e := range data {
		tmp.Status[i] = e
	}
	tmp.Mu.Unlock()
	index.Mu.Unlock()

	return nil
}

// CrawlRemoteRegistry scrapes all nicknames and user URLs
// from a provided registry. The urlKey passed to this function
// must be in the form of https://registry.example.com/api/plain/users
func (index *Index) CrawlRemoteRegistry(urlKey string) error {
	//fetch the remote registry's entries
	out, registry, err := GetTwtxt(urlKey)
	if err != nil {
		return err
	}

	// if we're working with an individual's twtxt file, error out
	if !registry {
		return fmt.Errorf("can't add single user via call to CrawlRemoteRegistry")
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
func (index *Index) GetUserStatuses(urlKey string) (TimeMap, error) {

	// Check that we have an initialized index.
	if index == nil {
		return nil, fmt.Errorf("can't get statuses from an empty index")

	} else if urlKey == "" {
		// Check that the request is valid.
		return nil, fmt.Errorf("can't retrieve statuses of blank user")

	} else if !strings.HasPrefix(urlKey, "http") {
		// Check that we were provided a URL
		return nil, fmt.Errorf("invalid URL: %v", urlKey)
	}

	// Check if the user is in the index. If they
	// aren't, we can't return their statuses.
	index.Mu.RLock()
	if _, ok := index.Reg[urlKey]; !ok {
		index.Mu.RUnlock()
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}

	// Pull the user's statuses from the index.
	index.Reg[urlKey].Mu.RLock()
	status := index.Reg[urlKey].Status
	index.Reg[urlKey].Mu.RUnlock()
	index.Mu.RUnlock()

	return status, nil
}

// GetStatuses returns a TimeMap containing all statuses
// from all users in the Index.
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
		v.Mu.RLock()
		if v.Status == nil || len(v.Status) == 0 {
			// Skip a user's statuses if the map is uninitialized or zero length
			v.Mu.RUnlock()
			continue
		}
		for a, b := range v.Status {
			if _, ok := v.Status[a]; ok {
				statuses[a] = b
			}
		}
		v.Mu.RUnlock()
	}
	index.Mu.RUnlock()

	return statuses, nil
}
