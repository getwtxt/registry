package registry

import (
	"fmt"
	"net/http"
	"testing"
)

func constructTwtxt() []byte {
	index := initTestEnv()
	var resp []byte
	// iterates through each mock user's mock statuses
	for _, v := range index {
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
	url     string
	wantErr bool
	local   bool
}{
	{
		url:     "http://localhost:8080/twtxt.txt",
		wantErr: false,
		local:   true,
	},
	{
		url:     "https://example3.com/twtxt.txt",
		wantErr: true,
		local:   false,
	},
	{
		url:     "https://example3.com",
		wantErr: true,
		local:   false,
	},
	{
		url:     "file://init_test.go",
		wantErr: true,
		local:   false,
	},
	{
		url:     "/etc/passwd",
		wantErr: true,
		local:   false,
	},
}

// Test the function that yoinks the /twtxt.txt file
// for a given user.
func Test_GetTwtxt(t *testing.T) {

	http.Handle("/twtxt.txt", http.HandlerFunc(twtxtHandler))
	go http.ListenAndServe(":8080", nil)

	for _, tt := range getTwtxtCases {
		t.Run(tt.url, func(t *testing.T) {
			if tt.local {
				t.Skipf("Local-only test. Skipping ... \n")
			}
			out, err := GetTwtxt(tt.url)
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
			b.Run(tt.url, func(b *testing.B) {
				GetTwtxt(tt.url)
			})
		}
	}
}
