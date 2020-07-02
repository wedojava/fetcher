package fetcher

import (
	"fmt"
	"strings"

	"github.com/wedojava/gears"
	"golang.org/x/net/html"
)

// soleTitle returns the text of the first non-empty title element
// in doc, and an error if there was not exactly one.
func soleTitleMutex(doc *html.Node) (title string, err error) {
	type bailout struct{}
	defer func() {
		switch p := recover(); p {
		case nil:
			// no panic
		case bailout{}:
			// "expected" panic
			err = fmt.Errorf("multiple title elements")
		default:
			panic(p)
		}
	}()
	// Bail out of recursion if we find more than one non-empty title.
	ForEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			if title != "" {
				panic(bailout{}) // multiple title elements
			}
			title = n.FirstChild.Data
		}
	}, nil)
	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}

func (p *Post) SetTitleMutex() error {
	title, err := soleTitleMutex(p.DOC)
	if err != nil {
		return err
	}
	ReplaceIllegalChar(&title)
	p.Title = strings.TrimSpace(title)
	switch p.Domain {
	case "www.boxun.com":
		if err = gears.ConvertToUtf8(&p.Title, "gbk", "utf8"); err != nil {
			return err
		}
	case "www.dwnews.com":
		p.Title = title[:strings.Index(title, "｜")]
	}
	return nil
}

func (p *Post) SetTitle() error {
	n := ElementsByTagName(p.DOC, "title")
	title := n[0].FirstChild.Data
	ReplaceIllegalChar(&title)
	p.Title = strings.TrimSpace(title)
	switch p.Domain {
	case "www.boxun.com":
		if err := gears.ConvertToUtf8(&p.Title, "gbk", "utf8"); err != nil {
			return err
		}
	case "www.dwnews.com":
		p.Title = title[:strings.Index(title, "｜")]
	}
	return nil
}
