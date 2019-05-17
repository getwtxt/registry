// Package registry implements functions and types that assist
// in the creation and management of a twtxt registry.
package registry // import "github.com/getwtxt/registry"

import (
	"bufio"
	"net/http"
	"os"
	"reflect"
	"testing"
)

var addUserCases = []struct {
	name      string
	nick      string
	url       string
	wantErr   bool
	localOnly bool
}{
	{
		name:      "Legitimate User (Local Only)",
		nick:      "testuser1",
		url:       "http://localhost:8080/twtxt.txt",
		wantErr:   false,
		localOnly: true,
	},
	{
		name:      "Unreachable twtxt File",
		nick:      "testuser2",
		url:       "https://example555555555.com/twtxt.txt",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "Empty Query",
		nick:      "",
		url:       "",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "Invalid URL",
		nick:      "foo",
		url:       "foobarringtons",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "Garbage Data",
		nick:      "",
		url:       "",
		wantErr:   true,
		localOnly: false,
	},
}

// Tests if we can successfully add a user to the index
func Test_UserIndex_AddUser(t *testing.T) {
	index := initTestEnv()
	if !addUserCases[0].localOnly {
		http.Handle("/twtxt.txt", http.HandlerFunc(twtxtHandler))
		go http.ListenAndServe(":8080", nil)
	}
	var buf = make([]byte, 256)
	// read random data into case 5
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	reader.Read(buf)
	addUserCases[4].nick = string(buf)
	addUserCases[4].url = string(buf)

	for n, tt := range addUserCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.localOnly {
				t.Skipf("Local-only test. Skipping ... ")
			}

			err := index.AddUser(tt.nick, tt.url)

			// only run some checks if we don't want an error
			if !tt.wantErr {
				if err != nil {
					t.Errorf("Got error: %v\n", err)
				}

				// make sure we have *something* in the index
				if reflect.ValueOf(index.Reg[tt.url]).IsNil() {
					t.Errorf("Failed to add user %v index.\n", tt.url)
				}

				// see if the nick in the index is the same
				// as the test case. verifies the URL and the nick
				// since the URL is used as the key
				data := index.Reg[tt.url]
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
			index.Reg[tt.url] = &Data{}
		}
	}
}

var delUserCases = []struct {
	name    string
	url     string
	wantErr bool
}{
	{
		name:    "Valid User",
		url:     "https://example.com/twtxt.txt",
		wantErr: false,
	},
	{
		name:    "Valid User",
		url:     "https://example3.com/twtxt.txt",
		wantErr: false,
	},
	{
		name:    "Already Deleted User",
		url:     "https://example3.com/twtxt.txt",
		wantErr: true,
	},
	{
		name:    "Empty Query",
		url:     "",
		wantErr: true,
	},
	{
		name:    "Garbage Data",
		url:     "",
		wantErr: true,
	},
}

// Tests if we can successfully delete a user from the index
func Test_UserIndex_DelUser(t *testing.T) {
	index := initTestEnv()
	var buf = make([]byte, 256)
	// read random data into case 5
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	reader.Read(buf)
	delUserCases[4].url = string(buf)

	for n, tt := range delUserCases {
		t.Run(tt.name, func(t *testing.T) {

			err := index.DelUser(tt.url)
			if !reflect.ValueOf(index.Reg[tt.url]).IsNil() {
				t.Errorf("Failed to delete user %v from index.\n", tt.url)
			}
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but did not receive. Case %v\n", n)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for case %v: %v\n", n, err)
			}
		})
	}
}
func Benchmark_UserIndex_DelUser(b *testing.B) {
	index := initTestEnv()

	data1 := &Data{
		Nick:    index.Reg[delUserCases[0].url].Nick,
		Date:    index.Reg[delUserCases[0].url].Date,
		APIdate: index.Reg[delUserCases[0].url].APIdate,
		Status:  index.Reg[delUserCases[0].url].Status,
	}

	data2 := &Data{
		Nick:    index.Reg[delUserCases[1].url].Nick,
		Date:    index.Reg[delUserCases[1].url].Date,
		APIdate: index.Reg[delUserCases[1].url].APIdate,
		Status:  index.Reg[delUserCases[1].url].Status,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range delUserCases {
			index.DelUser(tt.url)
		}

		index.Reg[delUserCases[0].url] = data1
		index.Reg[delUserCases[1].url] = data2
	}
}

var getUserStatusCases = []struct {
	name    string
	url     string
	wantErr bool
}{
	{
		name:    "Valid User",
		url:     "https://example.com/twtxt.txt",
		wantErr: false,
	},
	{
		name:    "Valid User",
		url:     "https://example3.com/twtxt.txt",
		wantErr: false,
	},
	{
		name:    "Nonexistent User",
		url:     "https://doesn't.exist/twtxt.txt",
		wantErr: true,
	},
	{
		name:    "Empty Query",
		url:     "",
		wantErr: true,
	},
	{
		name:    "Garbage Data",
		url:     "",
		wantErr: true,
	},
}

// Checks if we can retrieve a single user's statuses
func Test_UserIndex_GetUserStatuses(t *testing.T) {
	index := initTestEnv()
	var buf = make([]byte, 256)
	// read random data into case 5
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	reader.Read(buf)
	getUserStatusCases[4].url = string(buf)

	for n, tt := range getUserStatusCases {
		t.Run(tt.name, func(t *testing.T) {

			statuses, err := index.GetUserStatuses(tt.url)

			if !tt.wantErr {
				if reflect.ValueOf(statuses).IsNil() {
					t.Errorf("Failed to pull statuses for user %v\n", tt.url)
				}
				// see if the function returns the same data
				// that we already have
				data := index.Reg[tt.url]
				if !reflect.DeepEqual(data.Status, statuses) {
					t.Errorf("Incorrect data retrieved as statuses for user %v.\n", tt.url)
				}
			}

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, received nil for case %v: %v\n", n, tt.url)
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

		statuses, err := index.GetStatuses()
		if reflect.ValueOf(statuses).IsNil() || err != nil {
			t.Errorf("Failed to pull all statuses. %v\n", err)
		}

		// Now do the same query manually to see
		// if we get the same result
		unionmap := NewTimeMap()
		for _, v := range index.Reg {
			for i, e := range v.Status {
				unionmap[i] = e
			}
		}
		if !reflect.DeepEqual(statuses, unionmap) {
			t.Errorf("Incorrect data retrieved as statuses.\n")
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
