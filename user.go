package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// AddUser inserts a new user into the Index.
func (index *Index) AddUser(nickname, urlKey string, rlen string, ipAddress net.IP, statuses TimeMap) error {

	if index == nil {
		return fmt.Errorf("can't add user to uninitialized index")

	} else if nickname == "" || urlKey == "" {
		return fmt.Errorf("both URL and Nick must be specified")

	} else if !strings.HasPrefix(urlKey, "http") {
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

	index.Mu.Lock()
	defer index.Mu.Unlock()

	if _, ok := index.Users[urlKey]; ok {
		return fmt.Errorf("user %v already exists", urlKey)
	}

	index.Users[urlKey] = &User{
		Mu:     sync.RWMutex{},
		Nick:   nickname,
		URL:    urlKey,
		RLen:   rlen,
		IP:     ipAddress,
		Date:   time.Now().Format(time.RFC3339),
		Status: statuses}

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
	defer index.Mu.RUnlock()

	if _, ok := index.Users[urlKey]; !ok {
		return nil, fmt.Errorf("provided url key doesn't exist in index")
	}

	index.Users[urlKey].Mu.RLock()
	userGot := index.Users[urlKey]
	index.Users[urlKey].Mu.RUnlock()

	return userGot, nil
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

	index.Mu.Lock()
	defer index.Mu.Unlock()

	if _, ok := index.Users[urlKey]; !ok {
		return fmt.Errorf("can't delete user %v, user doesn't exist", urlKey)
	}

	delete(index.Users, urlKey)

	return nil
}

// UpdateUser scrapes an existing user's remote twtxt.txt
// file. Any new statuses are added to the user's entry
// in the Index. If the remote twtxt data's reported
// Content-Length does not differ from what is stored,
// an error is returned.
func (index *Index) UpdateUser(urlKey string) error {
	if urlKey == "" || !strings.HasPrefix(urlKey, "http") {
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

	diff, err := index.DiffTwtxt(urlKey)
	if err != nil {
		return err
	} else if !diff {
		return fmt.Errorf("no new statuses available for %v", urlKey)
	}

	out, registry, err := GetTwtxt(urlKey)
	if err != nil {
		return err
	}

	if registry {
		return fmt.Errorf("attempting to update registry URL - users should be updated individually")
	}

	index.Mu.Lock()
	defer index.Mu.Unlock()
	user := index.Users[urlKey]

	user.Mu.Lock()
	defer user.Mu.Unlock()
	nick := user.Nick

	data, err := ParseUserTwtxt(out, nick, urlKey)
	if err != nil {
		return err
	}

	for i, e := range data {
		user.Status[i] = e
	}

	index.Users[urlKey] = user

	return nil
}

// CrawlRemoteRegistry scrapes all nicknames and user URLs
// from a provided registry. The urlKey passed to this function
// must be in the form of https://registry.example.com/api/plain/users
func (index *Index) CrawlRemoteRegistry(urlKey string) error {
	if urlKey == "" || !strings.HasPrefix(urlKey, "http") {
		return fmt.Errorf("invalid URL: %v", urlKey)
	}

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
	defer index.Mu.Unlock()
	for _, e := range users {
		if _, ok := index.Users[e.URL]; !ok {
			index.Users[e.URL] = e
		}
	}

	return nil
}

// GetUserStatuses returns a TimeMap containing single user's statuses
func (index *Index) GetUserStatuses(urlKey string) (TimeMap, error) {
	if index == nil {
		return nil, fmt.Errorf("can't get statuses from an empty index")
	} else if urlKey == "" || !strings.HasPrefix(urlKey, "http") {
		return nil, fmt.Errorf("invalid URL: %v", urlKey)
	}

	index.Mu.RLock()
	defer index.Mu.RUnlock()
	if _, ok := index.Users[urlKey]; !ok {
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}

	index.Users[urlKey].Mu.RLock()
	status := index.Users[urlKey].Status
	index.Users[urlKey].Mu.RUnlock()

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
	defer index.Mu.RUnlock()

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

	return statuses, nil
}
