package registry

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func constructTwtxt() []byte {
	index := initTestEnv()
	var resp []byte
	// iterates through each mock user's mock statuses
	for _, v := range index.Reg {
		for _, e := range v.Status {
			resp = append(resp, []byte(fmt.Sprintf(e+"\n"))...)
		}
	}
	return resp
}

// this is just dumping all the mock statuses.
// it'll be served under fake paths as
// "remote" twtxt.txt files
func twtxtHandler(w http.ResponseWriter, _ *http.Request) {
	// prepare the response
	resp := constructTwtxt()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	n, err := w.Write(resp)
	if err != nil || n == 0 {
		fmt.Printf("Got error or wrote zero bytes: %v bytes, %v\n", n, err)
	}
}

var getTwtxtCases = []struct {
	name      string
	url       string
	wantErr   bool
	localOnly bool
}{
	{
		name:      "http://localhost:8080/twtxt.txt",
		url:       "http://localhost:8080/twtxt.txt",
		wantErr:   false,
		localOnly: true,
	},
	{
		name:      "https://example33333333333.com/twtxt.txt",
		url:       "https://example33333333333.com/twtxt.txt",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "https://example333333333333.com",
		url:       "https://example333333333333.com",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "file://init_test.go",
		url:       "file://init_test.go",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "/etc/passwd",
		url:       "/etc/passwd",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "https://example.com/file.cgi",
		url:       "https://example.com/file.cgi",
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "Garbage Data",
		url:       "this will be replaced with garbage data",
		wantErr:   true,
		localOnly: true,
	},
}

// Test the function that yoinks the /twtxt.txt file
// for a given user.
func Test_GetTwtxt(t *testing.T) {
	var buf = make([]byte, 256)
	// read random data into case 4
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	n, err := reader.Read(buf)
	if err != nil || n == 0 {
		t.Errorf("Couldn't set up test: %v\n", err)
	}
	getTwtxtCases[6].url = string(buf)

	if !getTwtxtCases[0].localOnly {
		http.Handle("/twtxt.txt", http.HandlerFunc(twtxtHandler))
		go fmt.Println(http.ListenAndServe(":8080", nil))
	}

	for _, tt := range getTwtxtCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.localOnly {
				t.Skipf("Local-only test. Skipping ... \n")
			}
			out, _, err := GetTwtxt(tt.url)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error: %v\n", tt.url)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v %v\n", tt.url, err)
			}
			if !tt.wantErr && out == nil {
				t.Errorf("Incorrect data received: %v\n", out)
			}
		})
	}

}

// running the benchmarks separately for each case
// as they have different properties (allocs, time)
func Benchmark_GetTwtxt(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, tt := range getTwtxtCases {
			_, _, err := GetTwtxt(tt.url)
			if err != nil {
				continue
			}
		}
	}
}

var parseTwtxtCases = []struct {
	name      string
	data      []byte
	wantErr   bool
	localOnly bool
}{
	{
		name:      "Constructed twtxt file",
		data:      constructTwtxt(),
		wantErr:   false,
		localOnly: false,
	},
	{
		name:      "Incorrectly formatted date",
		data:      []byte("foo_barrington\thttps://example3.com/twtxt.txt\t2019 April 23rd\tI love twtxt!!!11"),
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "No data",
		data:      []byte{},
		wantErr:   true,
		localOnly: false,
	},
	{
		name:      "Random/garbage data",
		wantErr:   true,
		localOnly: true,
	},
}

// See if we can break ParseTwtxt or get it
// to throw an unexpected error
func Test_ParseTwtxt(t *testing.T) {
	var buf = make([]byte, 256)
	// read random data into case 4
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	n, err := reader.Read(buf)
	if err != nil || n == 0 {
		t.Errorf("Couldn't set up test: %v\n", err)
	}
	parseTwtxtCases[3].data = buf

	for _, tt := range parseTwtxtCases {
		if tt.localOnly {
			t.Skipf("Local-only test: Skipping ... \n")
		}
		t.Run(tt.name, func(t *testing.T) {

			timemap, errs := ParseTwtxt(tt.data)
			if errs == nil && tt.wantErr {
				t.Errorf("Expected error(s), received none.\n")
			}

			if !tt.wantErr {
				if errs != nil {
					t.Errorf("Unexpected error: %v\n", errs)
				}

				for k, v := range timemap {
					if k == (time.Time{}) || v == "" {
						t.Errorf("Empty status or empty timestamp: %v, %v\n", k, v)
					}
				}
			}
		})
	}
}

func Benchmark_ParseTwtxt(b *testing.B) {
	var buf = make([]byte, 256)
	// read random data into case 4
	rando, _ := os.Open("/dev/random")
	reader := bufio.NewReader(rando)
	n, err := reader.Read(buf)
	if err != nil || n == 0 {
		b.Errorf("Couldn't set up benchmark: %v\n", err)
	}
	parseTwtxtCases[3].data = buf

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range parseTwtxtCases {
			ParseTwtxt(tt.data)
		}
	}
}
