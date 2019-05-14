package registry // import "github.com/getwtxt/registry"

import (
	"reflect"
	"testing"
)

func Test_UserIndex_AddUser(t *testing.T) {
	index := initTestEnv()
	var addUserCases = []struct {
		nick string
		url  string
	}{
		{
			nick: "testuser1",
			url:  "https://example4.com/twtxt.txt",
		},
		{
			nick: "testuser2",
			url:  "https://example5.com/twtxt.txt",
		},
	}

	for _, tt := range addUserCases {
		t.Run(tt.nick, func(t *testing.T) {
			index.AddUser(tt.nick, tt.url)
			if reflect.ValueOf(index[tt.url]).IsNil() {
				t.Errorf("Failed to add user %v index.\n", tt.url)
			}
			data := index[tt.url]
			if data.Nick != tt.nick {
				t.Errorf("Incorrect user data added to index for user %v.\n", tt.url)
			}
		})
	}
}

func Test_UserIndex_DelUser(t *testing.T) {
	index := initTestEnv()
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
	for _, tt := range delUserCases {
		t.Run(tt.url, func(t *testing.T) {
			index.DelUser(tt.url)
			if !reflect.ValueOf(index[tt.url]).IsNil() {
				t.Errorf("Failed to delete user %v from index.\n", tt.url)
			}
		})
	}
}

func Test_UserIndex_GetUserStatuses(t *testing.T) {
	index := initTestEnv()
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
	for _, tt := range getUserStatusCases {
		t.Run(tt.url, func(t *testing.T) {
			statuses := index.GetUserStatuses(tt.url)
			if reflect.ValueOf(statuses).IsNil() {
				t.Errorf("Failed to pull statuses for user %v\n", tt.url)
			}
			data := index[tt.url]
			if !reflect.DeepEqual(data.Status, statuses) {
				t.Errorf("Incorrect data retrieved as statuses for user %v.\n", tt.url)
			}
		})
	}

}

func Test_UserIndex_GetStatuses(t *testing.T) {
	index := initTestEnv()
	t.Run("UserIndex.GetStatuses()", func(t *testing.T) {
		statuses := index.GetStatuses()
		if reflect.ValueOf(statuses).IsNil() {
			t.Errorf("Failed to pull all statuses.")
		}
	})
}
