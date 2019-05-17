// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
package registry // import "github.com/getwtxt/registry"

import (
	"sync"
	"time"
)

// Data on each user. Data.Nick is the specified nickname.
// Data.Date is the time.Time of the user's submission to
// the registry. Data.APIdate is the RFC3339-formatted
// date/time of the user's submission. Data.Status is a
// TimeMap containing the user's statuses.
type Data struct {
	Mu      sync.RWMutex
	Nick    string
	Date    time.Time
	APIdate []byte
	Status  TimeMap
}

// Index provides an index of users constructed from a
// map[string]*Data. A sync.RWMutex is included to restrict
// concurrent access to the map.
type Index struct {
	Mu  sync.RWMutex
	Reg map[string]*Data
}

// TimeMap holds extracted and processed user data as a
// string. A standard time.Time value is used as the key.
// The time.Time value is processed from the status's
// RFC3339 timestamp.
type TimeMap map[time.Time]string

// TimeMapSlice is a slice of TimeMap. Used for sorting the
// output of aggregate queries such as GetStatuses() by
// timestamp. Also used for combining the output of those
// same queries into a single TimeMap.
type TimeMapSlice []TimeMap

// TimeSlice is a slice of time.Time used for sorting
// a TimeMap by timestamp.
type TimeSlice []time.Time

// NewUserIndex returns an initialized Index and its
// associated sync.RWMutex
func NewUserIndex() *Index {
	return &Index{
		Mu:  sync.RWMutex{},
		Reg: make(map[string]*Data),
	}
}

// NewTimeMap returns an initialized TimeMap.
func NewTimeMap() TimeMap {
	return make(TimeMap)
}

// NewTimeMapSlice returns an initialized slice of
// TimeMaps with zero length.
func NewTimeMapSlice() TimeMapSlice {
	return make(TimeMapSlice, 0)
}

// Len returns the length of the TimeSlice to be sorted.
// This helps satisfy sort.Interface.
func (t TimeSlice) Len() int {
	return len(t)
}

// Less returns true if the timestamp at index i is before
// the timestamp at index j in a given TimeSlice.
// This helps satisfy sort.Interface.
func (t TimeSlice) Less(i, j int) bool {
	return t[i].Before(t[j])
}

// Swap transposes the timestamps at the two given indices
// for the TimeSlice receiver.
// This helps satisfy sort.Interface.
func (t TimeSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
