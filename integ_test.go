package registry

import (
	"flag"
	"strings"
	"testing"
)

var integration = flag.Bool("integration", false, "Signal to perform an integration test")

func Test_Integration(t *testing.T) {
	if !*integration {
		t.Skipf("Skipping integration test. Use `go test -v -args -integration` to perform.\n")
	}
	t.Logf("Creating index object ...\n")
	index := NewIndex()

	t.Logf("Fetching remote twtxt file ...\n")
	mainregistry, err := GetTwtxt("https://enotty.dk/soltempore.txt")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	t.Logf("Parsing remote twtxt file ...\n")
	parsed, errz := ParseTwtxt(mainregistry)
	if errz != nil {
		t.Errorf("%v\n", errz)
	}

	t.Logf("Adding new user to index ...\n")
	err = index.AddUser("TestRegistry", "https://enotty.dk/soltempore.txt", nil, parsed)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	t.Logf("Querying user statuses ...\n")
	queryuser, err := index.QueryUser("TestRegistry")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	for _, e := range queryuser {
		if !strings.Contains(e, "TestRegistry") {
			t.Errorf("QueryUser() returned incorrect data\n")
		}
	}

	t.Logf("Querying for keyword in statuses ...\n")
	querystatus, err := index.QueryInStatus("morning")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	for _, e := range querystatus {
		if !strings.Contains(e, "morning") {
			t.Errorf("QueryInStatus() returned incorrect data\n")
		}
	}

	t.Logf("Querying for all statuses ...\n")
	allstatus, err := index.QueryAllStatuses()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(allstatus) == 0 || allstatus == nil {
		t.Errorf("Got nil/zero from QueryAllStatuses")
	}

	t.Logf("Querying for all users ...\n")
	allusers, err := index.QueryUser("")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if len(allusers) == 0 || allusers == nil {
		t.Errorf("Got nil/zero users on empty QueryUser() query")
	}

	t.Logf("Deleting user ...\n")
	err = index.DelUser("https://enotty.dk/soltempore.txt")
	if err != nil {
		t.Errorf("%v\n", err)
	}
}
