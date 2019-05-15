package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"log"
	"time"
)

// AddUser inserts a new user into the index. The *Data struct
// contains the nickname and the time the user was added.
func (index UserIndex) AddUser(nick string, url string) error {

	if nick == "" || url == "" {
		return fmt.Errorf("both URL and Nick must be specified")
	}

	imutex.RLock()
	if _, ok := index[url]; ok {
		imutex.RUnlock()
		log.Printf("User %v can't be added - already exists.\n", url)
		return fmt.Errorf("user %v already exists", url)
	}
	imutex.RUnlock()

	thetime := time.Now()
	rfc3339date, err := thetime.MarshalText()
	if err != nil {
		log.Printf("Error formatting user add time as RFC3339: %v\n", err)
	}

	imutex.Lock()
	index[url] = &Data{Nick: nick, Date: time.Now(), APIdate: rfc3339date}
	imutex.Unlock()

	return nil
}

// DelUser removes a user from the index completely.
func (index UserIndex) DelUser(url string) error {

	if url == "" {
		return fmt.Errorf("can't delete blank user")
	}

	imutex.RLock()
	if _, ok := index[url]; !ok {
		imutex.RUnlock()
		return fmt.Errorf("can't delete user %v, user doesn't exist", url)
	}
	imutex.RUnlock()

	imutex.Lock()
	delete(index, url)
	imutex.Unlock()

	return nil
}

// GetUserStatuses returns a TimeMap containing a user's statuses
func (index UserIndex) GetUserStatuses(url string) (TimeMap, error) {

	if url == "" {
		return nil, fmt.Errorf("can't retrieve statuses of blank user")
	}

	imutex.RLock()
	if _, ok := index[url]; !ok {
		imutex.RUnlock()
		return nil, fmt.Errorf("can't retrieve statuses of nonexistent user")
	}
	imutex.RUnlock()

	imutex.RLock()
	status := index[url].Status
	imutex.RUnlock()

	return status, nil
}

// GetStatuses returns a TimeMap containing all statuses
func (index UserIndex) GetStatuses() (TimeMap, error) {

	if index == nil {
		return nil, fmt.Errorf("can't retrieve statuses from empty index")
	}

	statuses := NewTimeMap()

	imutex.RLock()
	for _, v := range index {
		for a, b := range v.Status {
			statuses[a] = b
		}
	}
	imutex.RUnlock()

	return statuses, nil
}
