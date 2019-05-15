package registry

import (
	"strings"
	"testing"
)

var queryUserCases = []struct {
	nick    string
	wantErr bool
}{
	{
		nick:    "foo",
		wantErr: false,
	},
	{
		nick:    "example",
		wantErr: true,
	},
	{
		nick:    "",
		wantErr: false,
	},
}

// Checks if UserIndex.QueryUser() returns users that
// match the provided substring.
func Test_UserIndex_QueryUser(t *testing.T) {
	index := initTestEnv()
	for n, tt := range queryUserCases {
		t.Run(tt.nick, func(t *testing.T) {
			out, err := index.QueryUser(tt.nick)
			if out == nil && err != nil && !tt.wantErr {
				t.Errorf("Received nil output or an error when unexpected. Case %v, %v, %v\n", n, tt.nick, err)
			}
			if out != nil && tt.wantErr {
				t.Errorf("Received unexpected nil output when an error was expected. Case %v, %v\n", n, tt.nick)
			}
			for _, e := range out {
				one := strings.Split(e, "\t")
				if !strings.Contains(one[0], tt.nick) {
					t.Errorf("Received incorrect output: %v != %v\n", tt.nick, e)
				}
			}
		})
	}
}
func Benchmark_UserIndex_QueryUser(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range queryUserCases {
			index.QueryUser(tt.nick)
		}
	}
}

var queryTagCases = []struct {
	tag     string
	wantNil bool
	wantErr bool
}{
	{
		tag:     "twtxt",
		wantNil: false,
		wantErr: false,
	},
	{
		tag:     "project",
		wantNil: false,
		wantErr: false,
	},
	{
		tag:     "",
		wantNil: true,
		wantErr: true,
	},
	{
		tag:     "ahfiurrenkhfkajdhfao",
		wantNil: true,
		wantErr: false,
	},
}

func Test_UserIndex_QueryTag(t *testing.T) {
	index := initTestEnv()
	for _, tt := range queryTagCases {
		t.Run(tt.tag, func(t *testing.T) {
			out, err := index.QueryTag(tt.tag)
			if err != nil && !tt.wantErr {
				t.Errorf("Caught unexpected error: %v\n", err)
			}
			if !tt.wantErr && out == nil && !tt.wantNil {
				t.Errorf("Got nil when expecting output\n")
			}
			if err == nil && tt.wantErr {
				t.Errorf("Expecting error, got nil.\n")
			}
			for _, e := range out {
				split := strings.Split(e, "\t")
				if !strings.Contains(split[3], tt.tag) {
					t.Errorf("Status without tag\n")
				}
			}
		})
	}

}
func Benchmark_UserIndex_QueryTag(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range queryTagCases {
			index.QueryTag(tt.tag)
		}
	}
}
func Test_Data_FindTag(t *testing.T) {
	index := initTestEnv()
	data := make([]*Data, 0)
	for _, v := range index {
		data = append(data, v)
	}
	i := 0
	for _, tt := range data {
		t.Run(tt.Nick, func(t *testing.T) {
			tag := tt.FindTag(queryTagCases[i].tag)
			if tag == nil {
				t.Errorf("Got nil tag\n")
			}
		})
	}

}
func Benchmark_Data_FindTag(b *testing.B) {
	index := initTestEnv()
	data := make([]*Data, 0)
	for _, v := range index {
		data = append(data, v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range data {
			for _, v := range queryTagCases {
				tt.FindTag(v.tag)
			}
		}
	}
}

//func Test_TimeMapSlice_SortByTime(t *testing.T) {

//}
//func Benchmark_TimeMapSlice_SortByTime(b *testing.B) {

//}
