package registry

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"
)

var queryUserCases = []struct {
	name    string
	term    string
	wantErr bool
}{
	{
		name:    "Valid User",
		term:    "foo",
		wantErr: false,
	},
	{
		name:    "Empty Query",
		term:    "",
		wantErr: false,
	},
	{
		name:    "Nonexistent User",
		term:    "doesntexist",
		wantErr: true,
	},
	{
		name:    "Garbage Data",
		term:    "will be replaced with garbage data",
		wantErr: true,
	},
}

// Checks if UserIndex.QueryUser() returns users that
// match the provided substring.
func Test_UserIndex_QueryUser(t *testing.T) {
	index := initTestEnv()
	var buf = make([]byte, 256)
	// read random data into case 8
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	n, err := reader.Read(buf)
	if err != nil || n == 0 {
		t.Errorf("Couldn't set up test: %v\n", err)
	}
	queryUserCases[3].term = string(buf)

	for n, tt := range queryUserCases {

		t.Run(tt.name, func(t *testing.T) {
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
			_, err := index.QueryUser(tt.term)
			if err != nil {
				continue
			}
		}
	}
}

var queryInStatusCases = []struct {
	name    string
	substr  string
	wantNil bool
	wantErr bool
}{
	{
		name:    "Tag in Status",
		substr:  "twtxt",
		wantNil: false,
		wantErr: false,
	},
	{
		name:    "Valid URL",
		substr:  "https://example.com/twtxt.txt",
		wantNil: false,
		wantErr: false,
	},
	{
		name:    "Multiple Words in Status",
		substr:  "next programming",
		wantNil: false,
		wantErr: false,
	},
	{
		name:    "Multiple Words, Not in Status",
		substr:  "explosive bananas from antarctica",
		wantNil: true,
		wantErr: false,
	},
	{
		name:    "Empty Query",
		substr:  "",
		wantNil: true,
		wantErr: true,
	},
	{
		name:    "Nonsense",
		substr:  "ahfiurrenkhfkajdhfao",
		wantNil: true,
		wantErr: false,
	},
	{
		name:    "Invalid URL",
		substr:  "https://doesnt.exist/twtxt.txt",
		wantNil: true,
		wantErr: false,
	},
	{
		name:    "Garbage Data",
		substr:  "will be replaced with garbage data",
		wantNil: true,
		wantErr: false,
	},
}

// This tests whether we can find a substring in all of
// the known status messages, disregarding the metadata
// stored with each status.
func Test_UserIndex_QueryInStatus(t *testing.T) {
	index := initTestEnv()
	var buf = make([]byte, 256)
	// read random data into case 8
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	n, err := reader.Read(buf)
	if err != nil || n == 0 {
		t.Errorf("Couldn't set up test: %v\n", err)
	}
	queryInStatusCases[7].substr = string(buf)

	for _, tt := range queryInStatusCases {

		t.Run(tt.name, func(t *testing.T) {

			out, err := index.QueryInStatus(tt.substr)
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
				split := strings.Split(string(e), "\t")

				if !strings.Contains(split[1], tt.substr) {
					t.Errorf("Status without substring returned\n")
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
			_, err := index.QueryInStatus(tt.substr)
			if err != nil {
				continue
			}
		}
	}
}

// Tests whether we can retrieve the 20 most
// recent statuses in the index
func Test_QueryLatestStatuses(t *testing.T) {
	index := initTestEnv()
	t.Run("Latest Statuses", func(t *testing.T) {
		out, err := index.QueryAllStatuses()
		if out == nil || len(out) > 20 || err != nil {
			t.Errorf("Got no statuses, or more than 20: %v, %v\n", len(out), err)
		}
	})
}
func Benchmark_QueryLatestStatuses(b *testing.B) {
	index := initTestEnv()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := index.QueryAllStatuses()
		if err != nil {
			continue
		}
	}
}

// This tests whether we can find a substring in the
// given user's status messages, disregarding the metadata
// stored with each status.
func Test_Data_FindInStatus(t *testing.T) {
	index := initTestEnv()
	var buf = make([]byte, 256)
	// read random data into case 8
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	n, err := reader.Read(buf)
	if err != nil || n == 0 {
		t.Errorf("Couldn't set up test: %v\n", err)
	}
	queryInStatusCases[7].substr = string(buf)

	data := make([]*Data, 0)

	for _, v := range index.Reg {
		data = append(data, v)
	}

	for _, tt := range queryInStatusCases {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range data {

				tag := e.FindInStatus(tt.substr)
				if tag == nil && !tt.wantNil {
					t.Errorf("Got nil tag\n")
				}
			}
		})
	}

}
func Benchmark_Data_FindInStatus(b *testing.B) {
	index := initTestEnv()
	data := make([]*Data, 0)

	for _, v := range index.Reg {
		data = append(data, v)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range data {
			for _, v := range queryInStatusCases {
				tt.FindInStatus(v.substr)
			}
		}
	}
}

func Test_TimeMapSlice_SortByTime(t *testing.T) {
	index := initTestEnv()

	statusmap, err := index.GetStatuses()
	if err != nil {
		t.Errorf("Failed to finish test initialization: %v\n", err)
	}

	statusmaps := NewTimeMapSlice()
	statusmaps = append(statusmaps, statusmap)

	t.Run("Sort By Time", func(t *testing.T) {
		sorted, err := statusmaps.SortByTime()
		if err != nil {
			t.Errorf("%v\n", err)
		}
		split := strings.Split(sorted[0], "\t")
		firsttime, _ := time.Parse("RFC3339", split[0])

		for i := range sorted {
			if i < len(sorted)-1 {

				nextsplit := strings.Split(sorted[i+1], "\t")
				nexttime, _ := time.Parse("RFC3339", nextsplit[0])

				if firsttime.Before(nexttime) {
					t.Errorf("Timestamps out of order: %v\n", sorted)
				}

				firsttime = nexttime
			}
		}
	})
}

// Benchmarking a sort of 1000000 statuses by timestamp.
// Right now it's at roughly 2000ns per 2 statuses.
// Set sortMultiplier to be the number of desired
// statuses divided by four.
func Benchmark_TimeMapSlice_SortByTime(b *testing.B) {
	// I set this to 250,000,000 and it hard-locked
	// my laptop. Oops.
	sortMultiplier := 250
	b.Logf("Benchmarking SortByTime with a constructed slice of %v statuses ...\n", sortMultiplier*4)
	index := initTestEnv()

	statusmap, err := index.GetStatuses()
	if err != nil {
		b.Errorf("Failed to finish benchmark initialization: %v\n", err)
	}

	// Constructed index has four statuses. This
	// makes a TimeMapSlice of 1000000 statuses.
	statusmaps := NewTimeMapSlice()
	for i := 0; i < sortMultiplier; i++ {
		statusmaps = append(statusmaps, statusmap)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := statusmaps.SortByTime()
		if err != nil {
			continue
		}
	}
}

func Test_TimeMap_SortByTime(t *testing.T) {
	index := initTestEnv()

	statusmap, err := index.GetStatuses()
	if err != nil {
		t.Errorf("Failed to finish test initialization: %v\n", err)
	}

	t.Run("Sort By Time", func(t *testing.T) {
		sorted, err := statusmap.SortByTime()
		if err != nil {
			t.Errorf("%v\n", err)
		}
		split := strings.Split(sorted[0], "\t")
		firsttime, _ := time.Parse("RFC3339", split[0])

		for i := range sorted {
			if i < len(sorted)-1 {

				nextsplit := strings.Split(sorted[i+1], "\t")
				nexttime, _ := time.Parse("RFC3339", nextsplit[0])

				if firsttime.Before(nexttime) {
					t.Errorf("Timestamps out of order: %v\n", sorted)
				}

				firsttime = nexttime
			}
		}
	})
}

func Benchmark_TimeMap_SortByTime(b *testing.B) {
	index := initTestEnv()

	statusmap, err := index.GetStatuses()
	if err != nil {
		b.Errorf("Failed to finish benchmark initialization: %v\n", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := statusmap.SortByTime()
		if err != nil {
			continue
		}
	}
}
