package rfa

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/wedojava/fetcher/internal/htmldoc"
	"github.com/wedojava/gears"
	"golang.org/x/net/html"
)

type Post struct {
	Domain   string
	URL      *url.URL
	DOC      *html.Node
	Raw      []byte
	Title    string
	Body     string
	Date     string
	Filename string
}

func SetPost(p *Post) error {
	if err := SetDate(p); err != nil {
		return err
	}
	if err := SetTitle(p); err != nil {
		return err
	}
	if err := SetBody(p); err != nil {
		return err
	}
	return nil
}

func SetDate(p *Post) error {
	doc := htmldoc.ElementsByTagAndType(p.DOC, "script", "application/ld+json")
	if doc == nil {
		return errors.New("[-] rfa SetDate err, cannot get target nodes.")
	}
	d := doc[0].FirstChild
	if d.Type != html.TextNode {
		return errors.New("[-] rfa SetDate err, target node have no text.")
	}
	raw := d.Data
	re := regexp.MustCompile(`"date\w*?":\s*?"(.*?)"`)
	rs := re.FindAllStringSubmatch(raw, -1)
	p.Date = rs[0][1] // dateModified -> rs[0][1], datePublished -> rs[1][1]
	return nil
}

func SetTitle(p *Post) error {
	n := htmldoc.ElementsByTagName(p.DOC, "title")
	title := n[0].FirstChild.Data
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return errors.New("[-] there is no DOC object to get and format.")
	}
	b, err := Rfa(p)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	h1 := fmt.Sprintf("# [%02d.%02d][%02d%02dH] %s", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	p.Body = fmt.Sprintf("%s\n\n%s", h1, b)
	return nil
}

func Rfa(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndId(doc, "div", "storytext")
	if len(nodes) == 0 {
		return "", errors.New(`[-] There is no element match '<div id="storytext">'`)
	}
	plist := htmldoc.ElementsByTagName(nodes[0], "p")
	for _, v := range plist {
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.Data == "b" {
			body += "** "
			blist := htmldoc.ElementsByTagName(v, "b")
			for _, b := range blist {
				_b := b.FirstChild
				if _b != nil && _b.Data != "" {
					body += b.FirstChild.Data
				}
			}
			body += " **  \n"
		} else {
			body += v.FirstChild.Data + "  \n"
		}
	}
	body = strings.ReplaceAll(body, "**   **  \n", "")
	body = strings.ReplaceAll(body, "br  \n", "")
	return body, nil
}
