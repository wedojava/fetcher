// Pacage links provides a link-extraction fuction.
package fetcher

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/wedojava/gears"
	"golang.org/x/net/html"
)

// ExtractLinks makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func ExtractLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		// TODO: compress layers
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				// append only the target website
				if strings.HasPrefix(a.Val, "http") && strings.Contains(a.Val, link.Hostname()) {
					links = append(links, link.String())
				} else if strings.HasPrefix(a.Val, "/") {
					links = append(links, link.String())
				}

			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func (f *Fetcher) SetLinks() error {
	_url, err := url.Parse(f.Entrance)
	if err != nil {
		return err
	}
	links, err := ExtractLinks(_url.String())
	if err != nil {
		log.Printf(`can't extract links from "%s": %s`, _url, err)
		return err
	}
	links = gears.StrSliceDeDupl(links)
	hostname := _url.Hostname()
	switch hostname {
	case "www.boxun.com":
		f.Links = LinksFilter(links, `.*?/news/.*/\d*.shtml`)
	case "www.dwnews.com":
		f.Links = LinksFilter(links, `.*?/.*?/\d{8}/`)
		kickOutLinksMatchPath(&f.Links, "zone")
		kickOutLinksMatchPath(&f.Links, "视觉")
	case "www.voachinese.com":
		f.Links = LinksFilter(links, `.*?/a/.*-.*.html`)
	case "www.rfa.org":
		f.Links = LinksFilter(links, `.*?/.*?-\d*.html`)
		kickOutLinksMatchPath(&f.Links, "about")
	}
	for i, l := range f.Links {
		fmt.Printf("%2d: %s\n", i+1, l)
	}
	return nil
}

// kickOutLinksMatchPath will kick out the links match the path,
// if path=="zone" it will kick out the links that contains "/zone/"
func kickOutLinksMatchPath(links *[]string, path string) {
	tmp := []string{}
	path = "/" + url.QueryEscape(path) + "/"
	for _, link := range *links {
		if !strings.Contains(link, path) {
			tmp = append(tmp, link)
		}
	}
	*links = tmp
}

// TODO: use point to impletement LinksFilter
// LinksFilter is support for SetLinks method
func LinksFilter(links []string, regex string) []string {
	flinks := []string{}
	re := regexp.MustCompile(regex)
	s := strings.Join(links, "\n")
	flinks = re.FindAllString(s, -1)
	return flinks
}
