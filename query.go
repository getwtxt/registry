package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"sort"
	"strings"
)

// QueryUser checks the user index for usernames or URLs that contain the
// term provided as an argument. Entries are returned sorted by the date
// they were added to the index. If the argument provided is blank, return
// all users.
func (index UserIndex) QueryUser(term string) ([]string, error) {

	if index == nil {
		return nil, fmt.Errorf("can't query empty index")
	}

	timekey := NewTimeMap()
	keys := make(TimeSlice, 0)
	var users []string

	imutex.RLock()
	for k, v := range index {
		if strings.Contains(v.Nick, term) || strings.Contains(k, term) {
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

// QueryInStatus returns all the known statuses that
// contain the provided substring (tag, mention URL, etc).
func (index UserIndex) QueryInStatus(substr string) ([]string, error) {

	if substr == "" {
		return nil, fmt.Errorf("cannot query for empty tag")
	}

	statusmap := NewTimeMapSlice()

	imutex.RLock()
	for _, v := range index {
		statusmap = append(statusmap, v.FindInStatus(substr))
	}
	imutex.RUnlock()

	return statusmap.SortByTime(), nil
}

// QueryLatestStatuses returns the 20 most recent statuses
// in the registry sorted by time.
func (index UserIndex) QueryLatestStatuses() ([]string, error) {
	statusmap, err := index.GetStatuses()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	statusmaps := NewTimeMapSlice()
	statusmaps = append(statusmaps, statusmap)
	sorted := statusmaps.SortByTime()

	if len(sorted) < 20 {
		return sorted, nil
	}
	return sorted[:19], nil
}

// FindInStatus takes a user's statuses and looks for a given substring.
// Returns the statuses with the substring as a []string.
func (userdata *Data) FindInStatus(word string) TimeMap {

	statuses := NewTimeMap()

	for k, e := range userdata.Status {
		parts := strings.Split(e, "\t")
		statusslice := strings.Split(parts[3], " ")

		for _, v := range statusslice {

			if strings.Contains(v, word) {
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
