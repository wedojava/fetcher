package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// Posts is the fetcher can return many results.
type Posts map[string]*ThePost

type ThePost struct {
	Entrance string
	Domain   string
	URL      string
	DOC      *html.Node
	Raw      string
	Title    string
	Body     string
	Date     string
}

type Paragraph struct {
	Type    string
	Content string
}

var originalHost string

// TODO: Can't use resp.Body twice, so Raw and DOC can't fetch at the sametime.
func (post *ThePost) SetRaw(retryTimeout time.Duration) error {
	posturl := post.URL
	// To judge if there is a syntex error on url
	url, err := url.Parse(posturl)
	if err != nil {
		return fmt.Errorf("bad url: %s", err)
	}
	if originalHost == "" {
		originalHost = url.Host
	}
	if originalHost != url.Host {
		return nil
	}
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := http.Get(posturl)
		if err == nil { // success
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			post.Raw = string(raw)
			return nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil
}

func (post *ThePost) SetDOC(retryTimeout time.Duration) error {
	posturl := post.URL
	// To judge if there is a syntex error on url
	url, err := url.Parse(posturl)
	if err != nil {
		return fmt.Errorf("bad url: %s", err)
	}
	if originalHost == "" {
		originalHost = url.Host
	}
	if originalHost != url.Host {
		return nil
	}
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := http.Get(posturl)
		if err == nil { // success
			defer resp.Body.Close()
			post.DOC, err = html.Parse(resp.Body)
			if err != nil {
				return err
			}
			return nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil
}

// WaitForServer attempts to contact the server of a URL.
// It tries for one minute using exponential back-off.
// It reports an error if all attemps fail.
func (post *ThePost) WaitForServer() error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(post.URL)
		if err == nil {
			return err // success
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to respond after %s", post.URL, timeout)
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		log.Printf(`can't extract links from "%s": %s`, url, err)
	}
	return list
}
