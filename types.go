// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
package registry // import "github.com/getwtxt/registry"

import (
	"net"
	"sync"
	"time"
)

// Data on each user. Used as an entry in
// Index.Reg with Data.URL as the key.
type Data struct {
	// Provided to aid in concurrency-safe
	// reads and writes. In most cases, the
	// "outer" RWMutex in the Index should be
	// used instead. This RWMutex is provided
	// should the library user need to access
	// a Data directly.
	Mu sync.RWMutex

	// Nick is the user-specified nickname.
	Nick string

	// The URL of the user's twtxt file
	URL string

	// The IP address of the user is optionally
	// recorded.
	IP net.IP

	// The timestamp, in RFC3339 format,
	// reflecting when the user was added.
	Date string

	// A TimeMap of the user's statuses
	// from their twtxt file.
	Status TimeMap
}

// Index provides an index of users. It holds the
// bulk of the registry data.
type Index struct {
	// Provided to aid in concurrency-safe
	// reads and writes to a given registry
	// index instance.
	Mu sync.RWMutex

	// The registry's user data is contained
	// in this map. The functions within this
	// library expect the key to be the URL of
	// a given user's twtxt file.
	Reg map[string]*Data
}

// TimeMap holds extracted and processed user data as a
// string. A standard time.Time value is used as the key.
// The time.Time value is processed from the status's
// RFC3339 timestamp.
type TimeMap map[time.Time]string

// TimeSlice is a slice of time.Time used for sorting
// a TimeMap by timestamp.
type TimeSlice []time.Time

// NewUserData returns a pointer to an initialized Data
func NewUserData() *Data {
	return &Data{
		Mu:     sync.RWMutex{},
		Status: NewTimeMap(),
	}
}

// NewIndex returns an initialized Index
func NewIndex() *Index {
	return &Index{
		Mu:  sync.RWMutex{},
		Reg: make(map[URLKey]*Data),
	}
}

// NewTimeMap returns an initialized TimeMap.
func NewTimeMap() TimeMap {
	return make(TimeMap)
}

// Len returns the length of the TimeSlice to be sorted.
// This helps satisfy sort.Interface.
func (t TimeSlice) Len() int {
	return len(t)
}

// Less returns true if the timestamp at index i is after
// the timestamp at index j in a given TimeSlice. This results
// in a descending (reversed) sort order for timestamps rather
// than ascending.
// This helps satisfy sort.Interface.
func (t TimeSlice) Less(i, j int) bool {
	return t[i].After(t[j])
}

// Swap transposes the timestamps at the two given indices
// for the TimeSlice receiver.
// This helps satisfy sort.Interface.
func (t TimeSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
