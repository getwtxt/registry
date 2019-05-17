// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
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
func (index *Index) QueryUser(term string) ([]string, error) {

	if index == nil {
		return nil, fmt.Errorf("can't query empty index for user")
	}

	timekey := NewTimeMap()
	keys := make(TimeSlice, 0)
	var users []string

	index.Mu.RLock()
	for k, v := range index.Reg {
		if index.Reg[k] == nil {
			// Skip the user if their entry is uninitialized
			continue
		}
		if strings.Contains(v.Nick, term) || strings.Contains(k, term) {
			timekey[v.Date] = v.Nick + "\t" + k + "\t" + string(v.APIdate) + "\n"
			keys = append(keys, v.Date)
		}
	}
	index.Mu.RUnlock()

	sort.Sort(keys)
	for _, e := range keys {
		users = append(users, timekey[e])
	}

	return users, nil
}

// QueryInStatus returns all the known statuses that
// contain the provided substring (tag, mention URL, etc).
func (index *Index) QueryInStatus(substr string) ([]string, error) {
	if substr == "" {
		return nil, fmt.Errorf("cannot query for empty tag")
	} else if index == nil {
		return nil, fmt.Errorf("can't query statuses of empty index")
	}

	statusmap := NewTimeMapSlice()

	index.Mu.RLock()
	for _, v := range index.Reg {
		statusmap = append(statusmap, v.FindInStatus(substr))
	}
	index.Mu.RUnlock()

	return statusmap.SortByTime(), nil
}

// QueryAllStatuses returns all statuses in the registry,
// sorted by timestamp. Output is a slice of bytes, ready
// to send to the requester.
func (index *Index) QueryAllStatuses() ([]string, error) {
	if index == nil {
		return nil, fmt.Errorf("can't get latest statuses from empty index")
	}

	statusmap, err := index.GetStatuses()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	statusmaps := NewTimeMapSlice()
	statusmaps = append(statusmaps, statusmap)
	sorted := statusmaps.SortByTime()

	return sorted, nil
}

// FindInStatus takes a user's statuses and looks for a given substring.
// Returns the statuses with the substring as a TimeMap.
func (userdata *Data) FindInStatus(word string) TimeMap {
	if userdata == nil {
		return nil
	} else if len(word) > 140 {
		return nil
	}

	statuses := NewTimeMap()

	userdata.Mu.RLock()
	for k, e := range userdata.Status {
		if _, ok := userdata.Status[k]; !ok {
			continue
		}

		parts := strings.Split(e, "\t")
		if strings.Contains(parts[3], word) {
			statuses[k] = e
		}

	}
	userdata.Mu.RUnlock()

	return statuses
}

// SortByTime returns a string slice of the query results,
// sorted by time.Time. The receiver is a TimeMapSlice. the
// results are returned as a []byte.
func (tm TimeMapSlice) SortByTime() []string {
	if tm == nil {
		return nil
	}

	var unionmap = NewTimeMap()
	var times = make(TimeSlice, 0)
	var data []string

	for _, e := range tm {
		for k, v := range e {
			if _, ok := e[k]; ok {
				unionmap[k] = v
			}
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
