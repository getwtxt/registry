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
	}{{
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
				t.Errorf("Failed to add user to index.\n")
			}
		})
	}
}
