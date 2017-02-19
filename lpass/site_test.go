package lpass

import "testing"

func TestQuerySites(t *testing.T) {
	sites, err := QuerySites("google.com")
	if err != nil {
		t.Fatal(err)
	}

	if len(sites) == 0 {
		t.Fatal("Should have gotten one site")
	}
}
