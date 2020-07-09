// Pacage links provides a link-extraction fuction.
package fetcher

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/wedojava/fetcher/internal/htmldoc"
	"github.com/wedojava/gears"
)

func (f *Fetcher) SetLinks() error {
	links, err := htmldoc.ExtractLinks(f.Entrance.String())
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
		KickOutLinksMatchPath(&f.Links, "voaweishi")
	case "www.rfa.org":
		f.Links = LinksFilter(links, `.*?/.*?-\d*.html`)
		KickOutLinksMatchPath(&f.Links, "about")
	}
	return nil
}

// kickOutLinksMatchPath will kick out the links match the path,
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
