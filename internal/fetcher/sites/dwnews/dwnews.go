package dwnews

import (
	"errors"
	"fmt"
	"net/url"
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
	metas := htmldoc.MetasByName(p.DOC, "parsely-pub-date")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	p.Date = cs[0]
	return nil
}

func SetTitle(p *Post) error {
	n := htmldoc.ElementsByTagName(p.DOC, "title")
	title := n[0].FirstChild.Data
	title = strings.TrimSpace(title)
	if strings.Contains(title, "｜") {
		p.Title = title[:strings.Index(title, "｜")]
	}
	if strings.Contains(title, "|") {
		p.Title = title[:strings.Index(title, "|")]
	}
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return errors.New("[-] there is no DOC object to get and format.")
	}
	b, err := Dwnews(p)
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

func Dwnews(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagName(doc, "article")
	if len(nodes) == 0 {
		return "", errors.New("[-] There is no tag named `<article>` from: " + p.URL.String())
	}
	articleDoc := nodes[0].FirstChild
	plist := htmldoc.ElementsByTagName(articleDoc, "p")
	if articleDoc.FirstChild.Data == "div" { // to fetch the summary block
		// body += fmt.Sprintf("\n > %s  \n", plist[0].FirstChild.Data) // redundant summary
		body += fmt.Sprintf("\n > ")
	}
	for _, v := range plist { // the last item is `推荐阅读：`
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.FirstChild != nil && v.FirstChild.Data == "strong" {
			body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.FirstChild.Data)
		} else {
			ok := true

			for _, a := range v.Parent.Attr {
				if a.Key == "class" {
					switch a.Val {
					// if it is a info for picture, igonre!
					case "sc-bdVaJa iHZvIS":
						ok = false
					// if it is a twitter content, ignore!
					case "twitter-tweet":
						ok = false
					}
				}
			}
			if ok {
				body += v.FirstChild.Data + "  \n"
			}
		}
	}
	body = strings.ReplaceAll(body, "strong", "")
	body = strings.ReplaceAll(body, "** 推荐阅读： **", "")
	return body, nil
}
