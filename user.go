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

	if index == nil {
		return fmt.Errorf("can't add user to uninitialized index")

	} else if nickname == "" || urlKey == "" {
		return fmt.Errorf("both URL and Nick must be specified")

	} else if !strings.HasPrefix(urlKey, "http") {
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

	index.Mu.RLock()
	if _, ok := index.Users[urlKey]; ok {
		index.Mu.RUnlock()
		return fmt.Errorf("user %v already exists", urlKey)
	}
	index.Mu.RUnlock()

	thetime := time.Now()

	index.Mu.Lock()
	index.Users[urlKey] = &User{
		Mu:     sync.RWMutex{},
		Nick:   nickname,
		URL:    urlKey,
		IP:     ipAddress,
		Date:   thetime.Format(time.RFC3339),
		Status: statuses}
	index.Mu.Unlock()

	return nil
}

// Put inserts a given User into an Index. The User
// being pushed need only have the URL field filled.
// All other fields may be empty.
// This can be destructive: an existing User in the
// Index will be overwritten if its User.URL is the
// same as the User.URL being pushed.
func (index *Index) Put(user *User) error {
	if user == nil {
		return fmt.Errorf("can't push nil data to index")
	}
	if index == nil || index.Users == nil {
		return fmt.Errorf("can't push data to index: index uninitialized")
	}
	user.Mu.RLock()
	if user.URL == "" {
		user.Mu.RUnlock()
		return fmt.Errorf("can't push data to index: missing URL for key")
	}
	urlKey := user.URL
	index.Mu.Lock()
	index.Users[urlKey] = user
	index.Mu.Unlock()
	user.Mu.RUnlock()

	return nil
}

// Get returns the User associated with the
// provided URL key in the Index.
func (index *Index) Get(urlKey string) (*User, error) {
	if index == nil {
		return nil, fmt.Errorf("can't pop from nil index")
	}
	if urlKey == "" {
		return nil, fmt.Errorf("can't pop unless provided a key")
	}

	index.Mu.RLock()
	if _, ok := index.Users[urlKey]; !ok {
		index.Mu.RUnlock()
		return nil, fmt.Errorf("provided url key doesn't exist in index")
	}

	index.Users[urlKey].Mu.RLock()
	userUser := index.Users[urlKey]
	index.Users[urlKey].Mu.RUnlock()
	index.Mu.RUnlock()

	return userUser, nil
}

// DelUser removes a user and all associated data from
// the Index.
func (index *Index) DelUser(urlKey string) error {

	if index == nil {
		return fmt.Errorf("can't delete user from empty index")

	} else if urlKey == "" {
		return fmt.Errorf("can't delete blank user")

	} else if !strings.HasPrefix(urlKey, "http") {
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

	index.Mu.RLock()
	if _, ok := index.Users[urlKey]; !ok {
		index.Mu.RUnlock()
		return fmt.Errorf("can't delete user %v, user doesn't exist", urlKey)
	}
	index.Mu.RUnlock()

	// The User mutex is never unlocked because
	// the User is deleted. It is only acquired to
	// prevent a panic if another thread is reading/writing
	// to the user.
	index.Mu.Lock()
	index.Users[urlKey].Mu.Lock()
	delete(index.Users, urlKey)
	index.Mu.Unlock()

	return nil
}

// UpdateUser scrapes an existing user's remote twtxt.txt
// file. Any new statuses are added to the user's entry
// in the Index.
func (index *Index) UpdateUser(urlKey string) error {
	out, registry, err := GetTwtxt(urlKey)
	if err != nil {
		return err
	}

	if registry {
		return fmt.Errorf("attempting to update registry URL - users should be updated individually")
	}

	index.Mu.RLock()
	user := index.Users[urlKey]
	index.Mu.RUnlock()

	user.Mu.RLock()
	nick := user.Nick
	user.Mu.RUnlock()

	data, err := ParseUserTwtxt(out, nick, urlKey)
	if err != nil {
		return err
	}
	user.Mu.Lock()
	for i, e := range data {
		user.Status[i] = e
	}
	user.Mu.Unlock()

	index.Mu.Lock()
	index.Users[urlKey] = user
	index.Mu.Unlock()

	return nil
}

// CrawlRemoteRegistry scrapes all nicknames and user URLs
// from a provided registry. The urlKey passed to this function
// must be in the form of https://registry.example.com/api/plain/users
func (index *Index) CrawlRemoteRegistry(urlKey string) error {
	out, registry, err := GetTwtxt(urlKey)
	if err != nil {
		return err
	}

	if !registry {
		return fmt.Errorf("can't add single user via call to CrawlRemoteRegistry")
	}

	users, err := ParseRegistryTwtxt(out)
	if err != nil {
		return err
	}

	// only add new users so we don't overwrite data
	// we already have (and lose statuses, etc)
	index.Mu.Lock()
	for _, e := range users {
		if _, ok := index.Users[e.URL]; !ok {
			index.Users[e.URL] = e
		}
	}
	index.Mu.Unlock()

	return nil
}

// GetUserStatuses returns a TimeMap containing single user's statuses
func (index *Index) GetUserStatuses(urlKey string) (TimeMap, error) {

	if index == nil {
		return nil, fmt.Errorf("can't get statuses from an empty index")

	} else if urlKey == "" {
		return nil, fmt.Errorf("can't retrieve statuses of blank user")

	} else if !strings.HasPrefix(urlKey, "http") {
		return nil, fmt.Errorf("invalid URL: %v", urlKey)
	}

	index.Mu.RLock()
	if _, ok := index.Users[urlKey]; !ok {
		index.Mu.RUnlock()
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}

	index.Users[urlKey].Mu.RLock()
	status := index.Users[urlKey].Status
	index.Users[urlKey].Mu.RUnlock()
	index.Mu.RUnlock()

	return status, nil
}

// GetStatuses returns a TimeMap containing all statuses
// from all users in the Index.
func (index *Index) GetStatuses() (TimeMap, error) {

	if index == nil {
		return nil, fmt.Errorf("can't get statuses from an empty index")
	}

	statuses := NewTimeMap()

	index.Mu.RLock()
	for _, v := range index.Users {
		v.Mu.RLock()
		if v.Status == nil || len(v.Status) == 0 {
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
