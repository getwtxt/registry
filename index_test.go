package registry // import "github.com/getwtxt/registry"

import (
	"reflect"
	"testing"
)

var addUserCases = []struct {
	nick    string
	url     string
	wantErr bool
}{
	{
		nick:    "testuser1",
		url:     "https://example4.com/twtxt.txt",
		wantErr: false,
	},
	{
		nick:    "testuser2",
		url:     "https://example5.com/twtxt.txt",
		wantErr: false,
	},
	{
		nick:    "testuser1",
		url:     "https://example4.com/twtxt.txt",
		wantErr: true,
	},
	{
		nick:    "",
		url:     "",
		wantErr: true,
	},
}

// Tests if we can successfully add a user to the index
func Test_UserIndex_AddUser(t *testing.T) {
	index := initTestEnv()

	for n, tt := range addUserCases {
		t.Run(tt.nick, func(t *testing.T) {

			err := index.AddUser(tt.nick, tt.url)

			// only run some checks if we don't want an error
			if !tt.wantErr {
				if err != nil {
					t.Errorf("Got error: %v\n", err)
				}

				// make sure we have *something* in the index
				if reflect.ValueOf(index[tt.url]).IsNil() {
					t.Errorf("Failed to add user %v index.\n", tt.url)
				}

				// see if the nick in the index is the same
				// as the test case. verifies the URL and the nick
				// since the URL is used as the key
				data := index[tt.url]
				if data.Nick != tt.nick {
					t.Errorf("Incorrect user data added to index for user %v.\n", tt.url)
				}
			}
			// check for the cases that should throw an error
			if tt.wantErr && err == nil {
				t.Errorf("Expected error for case %v, got nil\n", n)
			}
		})
	}
}
func Benchmark_UserIndex_AddUser(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range addUserCases {
			index.AddUser(tt.nick, tt.url)
			index[tt.url] = &Data{}
		}
	}
}

var delUserCases = []struct {
	url string
}{
	{
		url: "https://example.com/twtxt.txt",
	},
	{
		url: "https://example3.com/twtxt.txt",
	},
}

// Tests if we can successfully delete a user from the index
func Test_UserIndex_DelUser(t *testing.T) {
	index := initTestEnv()

	for _, tt := range delUserCases {
		t.Run(tt.url, func(t *testing.T) {

			index.DelUser(tt.url)
			if !reflect.ValueOf(index[tt.url]).IsNil() {
				t.Errorf("Failed to delete user %v from index.\n", tt.url)
			}
		})
	}
}
func Benchmark_UserIndex_DelUser(b *testing.B) {
	index := initTestEnv()
	data1 := &Data{
		Nick:    index[delUserCases[0].url].Nick,
		Date:    index[delUserCases[0].url].Date,
		APIdate: index[delUserCases[0].url].APIdate,
		Status:  index[delUserCases[0].url].Status,
	}
	data2 := &Data{
		Nick:    index[delUserCases[1].url].Nick,
		Date:    index[delUserCases[1].url].Date,
		APIdate: index[delUserCases[1].url].APIdate,
		Status:  index[delUserCases[1].url].Status,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range delUserCases {
			index.DelUser(tt.url)
		}
		index[delUserCases[0].url] = data1
		index[delUserCases[1].url] = data2
	}
}

var getUserStatusCases = []struct {
	url string
}{
	{
		url: "https://example.com/twtxt.txt",
	},
	{
		url: "https://example3.com/twtxt.txt",
	},
}

// Checks if we can retrieve a single user's statuses
func Test_UserIndex_GetUserStatuses(t *testing.T) {
	index := initTestEnv()

	for _, tt := range getUserStatusCases {
		t.Run(tt.url, func(t *testing.T) {

			statuses := index.GetUserStatuses(tt.url)
			if reflect.ValueOf(statuses).IsNil() {
				t.Errorf("Failed to pull statuses for user %v\n", tt.url)
			}

			// see if the function returns the same data
			// that we already have
			data := index[tt.url]
			if !reflect.DeepEqual(data.Status, statuses) {
				t.Errorf("Incorrect data retrieved as statuses for user %v.\n", tt.url)
			}
		})
	}
}
func Benchmark_UserIndex_GetUserStatuses(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range getUserStatusCases {
			index.GetUserStatuses(tt.url)
		}
	}
}

// Tests if we can retrieve all user statuses at once
func Test_UserIndex_GetStatuses(t *testing.T) {
	index := initTestEnv()
	t.Run("UserIndex.GetStatuses()", func(t *testing.T) {

		statuses := index.GetStatuses()
		if reflect.ValueOf(statuses).IsNil() {
			t.Errorf("Failed to pull all statuses.")
		}

		// Now do the same query manually to see
		// if we get the same result
		unionmap := NewTimeMap()
		for _, v := range index {
			for i, e := range v.Status {
				unionmap[i] = e
			}
		}
		if !reflect.DeepEqual(statuses, unionmap) {
			t.Errorf("Incorrect data retrieved as statuses.")
		}
	})
}
func Benchmark_UserIndex_GetStatuses(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index.GetStatuses()
	}
}
