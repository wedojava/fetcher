package fetcher

import (
	"testing"
)

func TestSetLinks(t *testing.T) {
	var f = &Fetcher{
		Entrance: "https://www.rfa.org/mandarin/",
		// Entrance: "https://www.voachinese.com",
		Links: nil,
		Posts: nil,
	}
	err := f.SetLinks()
	if err != nil {
		t.Errorf("SetLinks fail!\n%s", err)
	}
}
