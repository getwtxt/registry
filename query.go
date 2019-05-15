package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"sort"
	"strings"
)

// QueryUser checks the user index for nicknames that contain the
// nickname provided as an argument. Entries are returned sorted
// by the date they were added to the index. If the argument
// provided is blank, return all users.
func (index UserIndex) QueryUser(name string) ([]string, error) {

	if index == nil {
		return nil, fmt.Errorf("can't query empty index")
	}

	timekey := NewTimeMap()
	keys := make(TimeSlice, 0)
	var users []string

	imutex.RLock()
	for k, v := range index {
		if strings.Contains(v.Nick, name) {
			timekey[v.Date] = v.Nick + "\t" + k + "\t" + string(v.APIdate) + "\n"
			keys = append(keys, v.Date)
		}
	}
	imutex.RUnlock()

	sort.Sort(keys)
	for _, e := range keys {
		users = append(users, timekey[e])
	}

	return users, nil
}

// QueryTag returns all the known statuses that
// contain the provided tag.
func (index UserIndex) QueryTag(tag string) ([]string, error) {

	if tag == "" {
		return nil, fmt.Errorf("cannot query for empty tag")
	}

	statusmap := NewTimeMapSlice()

	imutex.RLock()
	for _, v := range index {
		statusmap = append(statusmap, v.FindTag(tag))
	}
	imutex.RUnlock()

	return statusmap.SortByTime(), nil
}

// FindTag takes a user's tweets and looks for a given tag.
// Returns the tweets with the tag as a []string.
func (userdata *Data) FindTag(tag string) TimeMap {

	statuses := NewTimeMap()

	for k, e := range userdata.Status {
		parts := strings.Split(e, "\t")
		statusslice := strings.Split(parts[3], " ")

		for _, v := range statusslice {

			if strings.Contains(v, tag) {
				statuses[k] = e
				break
			}
		}
	}

	return statuses
}

// SortByTime returns a string slice of the query results
// sorted by time.Time
func (tm TimeMapSlice) SortByTime() []string {

	var unionmap = NewTimeMap()
	var times = make(TimeSlice, 0)
	var data []string

	for _, e := range tm {
		for k, v := range e {
			unionmap[k] = v
		}
	}

	for k := range unionmap {
		times = append(times, k)
	}

	sort.Sort(times)

	for _, e := range times {
		data = append(data, unionmap[e])
	}

	return data
}
