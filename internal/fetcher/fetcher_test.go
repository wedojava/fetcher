package fetcher

import (
	"testing"
)

// func TestSetLinks(t *testing.T) {
//         var f = &Fetcher{
//                 Entrance: "https://www.rfa.org/mandarin/",
//                 // Entrance: "https://www.voachinese.com",
//                 Links:    nil,
//                 LinksNew: nil,
//                 LinksOld: nil,
//         }
//         err := f.SetLinks()
//         if err != nil {
//                 t.Errorf("SetLinks fail!\n%s", err)
//         }
// }

func TestCrawl(t *testing.T) {
	breadthFirst(crawl, []string{"https://www.rfa.org/mandarin/"})
}
