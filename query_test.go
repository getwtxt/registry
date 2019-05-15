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
