// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
package registry // import "github.com/getwtxt/registry"

import (
	"net"
	"sync"
	"time"
)

// Data on each user. Used as an entry in the
// Index wrapper struct's map.
type Data struct {
	// Provided to aid in concurrency-safe
	// reads and writes. In most cases, the
	// "outer" RWMutex in the Index struct
	// should be used instead. This RWMutex
	// is provided should the library user need
	// to access a Data struct concurrently.
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

// TimeMapSlice is a slice of TimeMap. Used for sorting the
// output of aggregate queries such as GetStatuses() by
// timestamp. Also used for combining the output of those
// same queries into a single TimeMap.
type TimeMapSlice []TimeMap

// TimeSlice is a slice of time.Time used for sorting
// a TimeMap by timestamp.
type TimeSlice []time.Time

// NewUserData returns a pointer to an initialized Data
// struct.
func NewUserData() *Data {
	return &Data{
		Mu:     sync.RWMutex{},
		Status: NewTimeMap(),
	}
}

// NewIndex returns an initialized Index and its
// associated sync.RWMutex
func NewIndex() *Index {
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

// Less returns true if the timestamp at index i is after
// the timestamp at index j in a given TimeSlice. Results
// in a descending sort order for times.
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
