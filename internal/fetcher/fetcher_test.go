package fetcher

import (
	"log"
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	for {
		BreadthFirst(Crawl, []string{
			// "https://www.boxun.com/rolling.shtml",
			"https://www.dwnews.com",
			// "https://www.voachinese.com",
			// "https://www.rfa.org/mandarin/",
		})

		log.Println("Sleep a sec ...")
		time.Sleep(5 * time.Minute) // only useful by goroutine
	}
}
