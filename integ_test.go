package registry

import (
	"strings"
	"testing"
)

// This tests all the operations on an index.
func Test_Integration(t *testing.T) {
	var integration = func(t *testing.T) {
		t.Logf("Creating index object ...\n")
		index := NewIndex()

		t.Logf("Fetching remote twtxt file ...\n")
		mainregistry, _, err := GetTwtxt("https://enotty.dk/soltempore.txt")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		t.Logf("Parsing remote twtxt file ...\n")
		parsed, errz := ParseTwtxt(mainregistry, false)
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
	var integration2 = func() {
		t.Logf("Creating index object ...\n")
		index := NewIndex()

		t.Logf("Fetching remote twtxt file ...\n")
		mainregistry, _, err := GetTwtxt("https://enotty.dk/soltempore.txt")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		t.Logf("Parsing remote twtxt file ...\n")
		parsed, errz := ParseTwtxt(mainregistry, false)
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
	t.Run("Integration Test", integration)
	if !testing.Short() {
		allocs := testing.AllocsPerRun(5, integration2)
		t.Logf("%v\n", allocs)
	}
}
