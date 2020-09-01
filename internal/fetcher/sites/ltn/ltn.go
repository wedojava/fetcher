package ltn

import (
	"bytes"
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
	if err := setDate(p); err != nil {
		return err
	}
	if err := setTitle(p); err != nil {
		return err
	}
	if err := setBody(p); err != nil {
		return err
	}
	return nil
}

func setDate(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	metas := htmldoc.MetasByProperty(p.DOC, "article:published_time")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("SetData got nothing.")
	}
	p.Date = cs[0]
	return nil
}

func setTitle(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	n := htmldoc.ElementsByTag(p.DOC, "title")
	if n == nil {
		return fmt.Errorf("[-] there is no element <title>")
	}
	title := n[0].FirstChild.Data
	title = strings.ReplaceAll(title, " - 自由時報電子報", "")
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func setBody(p *Post) error {
	if p.DOC == nil {
		return fmt.Errorf("[-] p.DOC is nil")
	}
	b, err := ltn(p)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	h1 := fmt.Sprintf("# [%02d.%02d][%02d%02dH] %s", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	p.Body = h1 + "\n\n" + b + "\n\n原地址：" + p.URL.String()
	return nil
}

func ltn(p *Post) (string, error) {
	if p.Raw == nil {
		return "", fmt.Errorf("[-] p.Raw is nil")
	}
	raw := p.Raw
	// Fetch content nodes
	r := htmldoc.DivWithAttr2(raw, "data-desc", "內容頁")
	ps := [][]byte{}
	b := bytes.Buffer{}
	re := regexp.MustCompile(`<p>(.*?)</p>`)
	for _, v := range re.FindAllSubmatch(r, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return "", fmt.Errorf("no <p> matched")
	}
	re = regexp.MustCompile(`<iframe.*?</iframe>`)
	for _, p := range ps {
		p = re.ReplaceAll(p, []byte(""))
		b.Write(p)
		b.Write([]byte("  \n"))
	}

	return b.String(), nil
}
