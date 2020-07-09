package fetcher

import (
	"fmt"
	"strings"

	"github.com/wedojava/fetcher/internal/htmldoc"
	"github.com/wedojava/gears"
	"golang.org/x/net/html"
)

// TODO: rm while stable
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
	htmldoc.ForEachNode(doc, func(n *html.Node) {
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

// TODO: rm while stable
func (p *Post) SetTitleMutex() error {
	title, err := soleTitleMutex(p.DOC)
	if err != nil {
		return err
	}
	switch p.Domain {
	case "www.boxun.com":
		if err = gears.ConvertToUtf8(&title, "gbk", "utf8"); err != nil {
			return err
		}
	case "www.dwnews.com":
		p.Title = title[:strings.Index(title, "｜")]
	}
	gears.ReplaceIllegalChar(&title)
	p.Title = strings.TrimSpace(title)
	return nil
}

// TODO: rm while TreatTitle pass test
func (p *Post) SetTitle() error {
	n := htmldoc.ElementsByTagName(p.DOC, "title")
	title := n[0].FirstChild.Data
	switch p.Domain {
	case "www.boxun.com":
		if err := gears.ConvertToUtf8(&title, "gbk", "utf8"); err != nil {
			return err
		}
	case "www.dwnews.com":
		p.Title = title[:strings.Index(title, "｜")]
	}
	gears.ReplaceIllegalChar(&title)
	p.Title = strings.TrimSpace(title)
	return nil
}

func (p *Post) TreatTitle(f func(*string) error) error {
	n := htmldoc.ElementsByTagName(p.DOC, "title")
	title := n[0].FirstChild.Data
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	if err := f(&p.Title); err != nil {
		return err
	}
	return nil
}
