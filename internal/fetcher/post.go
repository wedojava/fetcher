package fetcher

import (
	"log"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

type Post struct {
	Domain string
	URL    *url.URL
	DOC    *html.Node
	Raw    []byte
	Title  string
	Body   string
	Date   string
}

type Paragraph struct {
	Type    string
	Content string
}

func PostFactory(rawurl string) *Post {
	url, err := url.Parse(rawurl)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Post{
		Domain: url.Hostname(),
		URL:    url,
	}
}

func (p *Post) SetPost() error {
	// set contents
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.Raw = raw
	p.DOC = doc
	// set Date
	if err := p.SetDate(); err != nil {
		return err
	}
	// set Title
	if err := p.SetTitle(); err != nil {
		return err
	}
	// set Body (get and format body)
	if err := p.SetBody(); err != nil {
		return err
	}
	return nil
}
