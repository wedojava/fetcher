package fetcher

import (
	"fmt"
	"path/filepath"

	"golang.org/x/net/html"
)

func soleDate(doc *html.Node, metaName string) (date string, err error) {
	type bailout struct{}
	defer func() {
		switch p := recover(); p {
		case nil:
			// no panic
		case bailout{}:
			// "expected" panic
			err = fmt.Errorf("multiple date elements")
		default:
			panic(p)
		}
	}()
	// Bail out of recursion if we find more than one non-empty date.
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			yes := false
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == metaName {
					yes = true
				}
				if yes && a.Key == "content" {
					date = a.Val
				}
			}
		}
	}, nil)
	if date == "" {
		return "", fmt.Errorf("no date element")
	}
	return
}

func (p *Post) DwnewsDateInMeta() error {
	date, err := soleDate(p.DOC, "parsely-pub-date")
	if err != nil {
		return err
	}
	p.Date = date
	return nil
}

func (p *Post) SetDate() error {
	switch p.Domain {
	case "www.boxun.com":
		a := filepath.Base(p.URL.Path)
		p.Date = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", a[:4], a[4:6], a[6:8], a[8:10], a[10:12], "00")
	case "www.dwnews.com":
		if err := p.DwnewsDateInMeta(); err != nil {
			return err
		}
	}
	return nil
}
