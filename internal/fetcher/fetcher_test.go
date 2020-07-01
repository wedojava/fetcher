package fetcher

import (
	"testing"
)

func TestSetLinks(t *testing.T) {
	var f = &Fetcher{
		Entrance: "https://www.rfa.org/mandarin/",
		// Entrance: "https://www.voachinese.com",
	}
	err := f.SetLinks()
	if err != nil {
		t.Errorf("SetLinks fail!\n%s", err)
	}
}

func TestCrawl(t *testing.T) {
	breadthFirst(crawl, []string{
		"https://www.boxun.com/rolling.shtml",
		// "https://www.dwnews.com",
		// "https://www.voachinese.com",
		// "https://www.rfa.org/mandarin/",
	})
}
