package fetcher

import (
	"net/url"
	"testing"
)

func TestSetLinks(t *testing.T) {
	u, err := url.Parse("https://www.rfa.org/mandarin/")
	if err != nil {
		t.Errorf("Url Parse fail!\n%s", err)
	}
	var f = &Fetcher{
		Entrance: u,
		// Entrance: "https://www.voachinese.com",
	}
	f.SetLinks()
}

func TestCrawl(t *testing.T) {
	breadthFirst(crawl, []string{
		// "https://www.boxun.com/rolling.shtml",
		// "https://www.dwnews.com",
		"https://www.voachinese.com",
		// "https://www.rfa.org/mandarin/",
	})
}
