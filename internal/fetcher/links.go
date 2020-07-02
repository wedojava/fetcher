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
func ExtractLinks(str string) ([]string, error) {
	resp, err := http.Get(str)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", str, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", str, err)
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
	ForEachNode(doc, visitNode, nil)
	return links, nil
}

func (f *Fetcher) SetLinks() error {
	links, err := ExtractLinks(f.Entrance.String())
	if err != nil {
		log.Printf(`can't extract links from "%s": %s`, f.Entrance.String(), err)
		return err
	}
	links = gears.StrSliceDeDupl(links)
	hostname := f.Entrance.Hostname()
	switch hostname {
	case "www.boxun.com":
		f.Links = LinksFilter(links, `.*?/news/.*/\d*.shtml`)
	case "www.dwnews.com":
		f.Links = LinksFilter(links, `.*?/.*?/\d{8}/`)
		KickOutLinksMatchPath(&f.Links, "zone")
		KickOutLinksMatchPath(&f.Links, "/"+url.QueryEscape("视觉")+"/")
	case "www.voachinese.com":
		f.Links = LinksFilter(links, `.*?/a/.*-.*.html`)
	case "www.rfa.org":
		f.Links = LinksFilter(links, `.*?/.*?-\d*.html`)
		KickOutLinksMatchPath(&f.Links, "about")
	}
	return nil
}

// kickOutLinksMatchPath will kick out the links match the path,
// if path=="zone" it will kick out the links that contains "/zone/"
func KickOutLinksMatchPath(links *[]string, path string) {
	tmp := []string{}
	// path = "/" + url.QueryEscape(path) + "/"
	// path = url.QueryEscape(path)
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
