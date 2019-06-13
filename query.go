package registry // import "github.com/getwtxt/registry"

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// QueryUser checks the Index for usernames
// or user URLs that contain the term provided as an argument. Entries
// are returned sorted by the date they were added to the Index. If
// the argument provided is blank, return all users.
func (index *Index) QueryUser(term string) ([]string, error) {
	if index == nil {
		return nil, fmt.Errorf("can't query empty index for user")
	}

	term = strings.ToLower(term)
	timekey := NewTimeMap()
	keys := make(TimeSlice, 0)
	var users []string

	index.Mu.RLock()
	defer index.Mu.RUnlock()

	for k, v := range index.Users {
		if index.Users[k] == nil {
			continue
		}
		v.Mu.RLock()
		if strings.Contains(strings.ToLower(v.Nick), term) || strings.Contains(strings.ToLower(k), term) {
			thetime, err := time.Parse(time.RFC3339, v.Date)
			if err != nil {
				v.Mu.RUnlock()
				continue
			}
			timekey[thetime] = v.Nick + "\t" + k + "\t" + v.Date + "\n"
			keys = append(keys, thetime)
		}
		v.Mu.RUnlock()
	}

	sort.Sort(keys)
	for _, e := range keys {
		users = append(users, timekey[e])
	}

	return users, nil
}

// QueryInStatus returns all statuses in the Index
// that contain the provided substring (tag, mention URL, etc).
func (index *Index) QueryInStatus(substring string) ([]string, error) {
	if substring == "" {
		return nil, fmt.Errorf("cannot query for empty tag")
	} else if index == nil {
		return nil, fmt.Errorf("can't query statuses of empty index")
	}

	statusmap := make([]TimeMap, 0)

	index.Mu.RLock()
	defer index.Mu.RUnlock()

	for _, v := range index.Users {
		statusmap = append(statusmap, v.FindInStatus(substring))
	}

	sorted, err := SortByTime(statusmap...)
	if err != nil {
		return nil, err
	}

	return sorted, nil
}

// QueryAllStatuses returns all statuses in the Index
// as a slice of strings sorted by timestamp.
func (index *Index) QueryAllStatuses() ([]string, error) {
	if index == nil {
		return nil, fmt.Errorf("can't get latest statuses from empty index")
	}

	statusmap, err := index.GetStatuses()
	if err != nil {
		return nil, err
	}

	sorted, err := SortByTime(statusmap)
	if err != nil {
		return nil, err
	}

	return sorted, nil
}

// ReduceToPage returns the passed 'page' worth of output.
// One page is twenty items. For example, if 2 is passed,
// it will return data[20:40]. According to the twtxt
// registry specification, queries should accept a "page"
// value.
func ReduceToPage(page int, data []string) []string {
	end := 20 * page
	if end > len(data) || end < 1 {
		end = len(data)
	}

	beg := end - 20
	if beg > len(data)-1 || beg < 0 {
		beg = 0
	}

	return data[beg:end]
}

// FindInStatus takes a user's statuses and looks for a given substring.
// Returns the statuses that include the substring as a TimeMap.
func (userdata *User) FindInStatus(substring string) TimeMap {
	if userdata == nil {
		return nil
	} else if len(substring) > 140 {
		return nil
	}

	substring = strings.ToLower(substring)
	statuses := NewTimeMap()

	userdata.Mu.RLock()
	defer userdata.Mu.RUnlock()

	for k, e := range userdata.Status {
		if _, ok := userdata.Status[k]; !ok {
			continue
		}

		parts := strings.Split(strings.ToLower(e), "\t")
		if strings.Contains(parts[3], substring) {
			statuses[k] = e
		}
	}

	return statuses
}

// SortByTime returns a string slice of the query results,
// sorted by timestamp in descending order (newest first).
func SortByTime(tm ...TimeMap) ([]string, error) {
	if tm == nil {
		return nil, fmt.Errorf("can't sort nil TimeMaps")
	}

	var times = make(TimeSlice, 0)
	var data []string

	for _, e := range tm {
		for k := range e {
			times = append(times, k)
		}
	}

	sort.Sort(times)

	for k := range tm {
		for _, e := range times {
			if _, ok := tm[k][e]; ok {
				data = append(data, tm[k][e])
			}
		}
	}

	return data, nil
}
