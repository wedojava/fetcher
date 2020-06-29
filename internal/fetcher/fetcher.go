package fetcher

import (
	"fmt"
	"github.com/wedojava/gears"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Fetcher struct {
	Entrance string
	Links    []string
	LinksNew []string
	LinksOld []string
}

var originalHost string

// TODO: Can't use resp.Body twice, so Raw and DOC can't fetch at the sametime.
// GetRaw can get html raw bytes by rawurl.
func GetRaw(url *url.URL, retryTimeout time.Duration) ([]byte, error) {
	// To judge if there is a syntex error on url
	if originalHost == "" {
		originalHost = url.Host
	}
	if originalHost != url.Host {
		return nil, fmt.Errorf("bad host of url: %s, expected: %s", url.Host, originalHost)
	}
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := http.Get(url.String())
		if err == nil { // success
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return raw, nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil, nil
}

func GetDOC(url *url.URL, retryTimeout time.Duration) (*html.Node, error) {
	// To judge if there is a syntex error on url
	if originalHost == "" {
		originalHost = url.Host
	}
	if originalHost != url.Host {
		return nil, fmt.Errorf("bad host of url: %s, expected: %s", url.Host, originalHost)
	}
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := http.Get(url.String())
		if err == nil { // success
			defer resp.Body.Close()
			doc, err := html.Parse(resp.Body)
			if err != nil {
				return nil, err
			}
			return doc, nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil, nil
}

func (f *Fetcher) SetLinks() error {
	url, err := url.Parse(f.Entrance)
	if err != nil {
		return err
	}
	links, err := ExtractLinks(url.String())
	if err != nil {
		log.Printf(`can't extract links from "%s": %s`, url, err)
		return err
	}
	links = gears.StrSliceDeDupl(links)
	hostname := url.Hostname()
	switch hostname {
	case "www.boxun.com":
		f.Links = LinksFilter(links, `.*?/news/.*/\d*.shtml`)
	case "www.dwnews.com":
		f.Links = LinksFilter(links, `.*?/.{2}/\d{8}/.+?`)
	case "www.voachinese.com":
		f.Links = LinksFilter(links, `.*?/a/.*-.*.html`)
	case "www.rfa.org":
		links = LinksFilter(links, `.*?/.*?-\d*.html`)
		for _, link := range links {
			if !strings.Contains(link, "/about/") {
				f.Links = append(f.Links, link)
			}
		}
	}
	// for i, l := range f.Links {
	//         fmt.Printf("%2d: %s\n", i+1, l)
	// }
	return nil
}

func LinksFilter(links []string, regex string) []string {
	flinks := []string{}
	re := regexp.MustCompile(regex)
	s := strings.Join(links, "\n")
	flinks = re.FindAllString(s, -1)
	return flinks
}

func FetcherFactory(site string) *Fetcher {
	return &Fetcher{
		Entrance: site,
		Links:    nil,
		LinksNew: nil,
		LinksOld: nil,
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func breadthFirst(f func(item string) error, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				f(item)
				worklist = items
				// worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) error {
	f := FetcherFactory(url)
	log.Printf("[*] Deal with: [%s]\n", url)
	log.Println("[*] Fetch links ...")
	if err := f.SetLinks(); err != nil {
		log.Println(err)
		return err
	}
	// Set LinksNew
	f.LinksNew = gears.StrSliceDiff(f.Links, f.LinksOld)
	// GetNews then compare via md5 and Save or Rewrite news exist
	for _, link := range f.LinksNew {
		post := PostFactory(link)
	}
	// Set LinksOld
	f.LinksOld = f.Links
	return nil
}
