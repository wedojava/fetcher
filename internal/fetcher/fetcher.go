package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/wedojava/gears"
	"golang.org/x/net/html"
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

func FetcherFactory(site string) *Fetcher {
	return &Fetcher{
		Entrance: site,
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func breadthFirst(f func(item string) error, worklist []string) {
	for _, item := range worklist {
		if err := f(item); err != nil {
			log.Println(err)
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
	log.Println("[*] Get news ...")
	for _, link := range f.LinksNew {
		post := PostFactory(link)
		if err := post.SetPost(); err != nil {
			return err
		}
	}
	// Set LinksOld
	f.LinksOld = f.Links
	return nil
}
