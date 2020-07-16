package main

import (
	"log"
	"strconv"
	"time"

	"github.com/wedojava/fetcher/internal/fetcher"
)

func main() {
	year := strconv.Itoa(time.Now().Year())
	sites := []string{
		"https://www.dwnews.com",
		"https://www.voachinese.com",
		"https://www.voachinese.com/z/1739",
		"https://www.rfa.org/mandarin/",
		"https://www.rfa.org/mandarin/Xinwen/story_archive?year=" + year,
		"https://www.rfa.org/mandarin/yataibaodao/story_archive?year=" + year,
		"https://www.boxun.com/rolling.shtml",
	}
	for {
		fetcher.BreadthFirst(fetcher.Crawl, sites)
		log.Println("Hold a sec ...")
		time.Sleep(5 * time.Minute)
	}
}
