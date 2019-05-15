package registry

import (
	"strings"
	"testing"
)

var queryUserCases = []struct {
	term    string
	wantErr bool
}{
	{
		term:    "foo",
		wantErr: false,
	},
	{
		term:    "example",
		wantErr: false,
	},
	{
		term:    "",
		wantErr: false,
	},
	{
		term:    "doesntexist",
		wantErr: true,
	},
}

// Checks if UserIndex.QueryUser() returns users that
// match the provided substring.
func Test_UserIndex_QueryUser(t *testing.T) {
	index := initTestEnv()

	for n, tt := range queryUserCases {

		t.Run(tt.term, func(t *testing.T) {
			out, err := index.QueryUser(tt.term)

			if out == nil && err != nil && !tt.wantErr {
				t.Errorf("Received nil output or an error when unexpected. Case %v, %v, %v\n", n, tt.term, err)
			}

			if out != nil && tt.wantErr {
				t.Errorf("Received unexpected nil output when an error was expected. Case %v, %v\n", n, tt.term)
			}

			for _, e := range out {
				one := strings.Split(e, "\t")

				if !strings.Contains(one[0], tt.term) && !strings.Contains(one[1], tt.term) {
					t.Errorf("Received incorrect output: %v != %v\n", tt.term, e)
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
			index.QueryUser(tt.term)
		}
	}
}

var queryInStatusCases = []struct {
	tag     string
	mention string
	wantNil bool
	wantErr bool
}{
	{
		tag:     "twtxt",
		mention: "https://example.com/twtxt.txt",
		wantNil: false,
		wantErr: false,
	},
	{
		tag:     "project",
		mention: "https://example3.com/twtxt.txt",
		wantNil: false,
		wantErr: false,
	},
	{
		tag:     "",
		mention: "",
		wantNil: true,
		wantErr: true,
	},
	{
		tag:     "ahfiurrenkhfkajdhfao",
		mention: "https://doesnt.exist/twtxt.txt",
		wantNil: true,
		wantErr: false,
	},
}

// This tests whether we can find a substring in all of
// the known status messages, disregarding the metadata
// stored with each status.
func Test_UserIndex_QueryInStatus(t *testing.T) {
	index := initTestEnv()

	for _, tt := range queryInStatusCases {

		t.Run(tt.tag, func(t *testing.T) {

			out, err := index.QueryInStatus(tt.tag)
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
			out, err = index.QueryInStatus(tt.mention)
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

				if !strings.Contains(split[3], tt.mention) {
					t.Errorf("Status without tag\n")
				}
			}
		})
	}

}
func Benchmark_UserIndex_QueryInStatus(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range queryInStatusCases {
			index.QueryInStatus(tt.tag)
		}
	}
}

// Tests whether we can retrieve the 20 most
// recent statuses in the index
func Test_QueryLatestStatuses(t *testing.T) {
	index := initTestEnv()
	t.Run("Latest Statuses", func(t *testing.T) {
		out, err := index.QueryLatestStatuses()
		if out == nil || len(out) > 20 || err != nil {
			t.Errorf("Got no statuses, or more than 20: %v, %v\n", len(out), err)
		}
	})
}
func Benchmark_QueryLatestStatuses(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index.QueryLatestStatuses()
	}
}

// This tests whether we can find a substring in the
// given user's status messages, disregarding the metadata
// stored with each status.
func Test_Data_FindInStatus(t *testing.T) {
	index := initTestEnv()
	data := make([]*Data, 0)

	for _, v := range index {
		data = append(data, v)
	}

	i := 0
	for _, tt := range data {
		t.Run(tt.Nick, func(t *testing.T) {

			tag := tt.FindInStatus(queryInStatusCases[i].tag)
			if tag == nil {
				t.Errorf("Got nil tag\n")
			}
			i++
		})
	}

}
func Benchmark_Data_FindInStatus(b *testing.B) {
	index := initTestEnv()
	data := make([]*Data, 0)

	for _, v := range index {
		data = append(data, v)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range data {
			for _, v := range queryInStatusCases {
				tt.FindInStatus(v.tag)
			}
		}
	}
}

//func Test_TimeMapSlice_SortByTime(t *testing.T) {

//}
//func Benchmark_TimeMapSlice_SortByTime(b *testing.B) {

//}
