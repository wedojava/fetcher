package rfa

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/wedojava/fetcher/internal/htmldoc"
)

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

func TestSetPost(t *testing.T) {
	p := PostFactory("https://www.rfa.org/mandarin/Xinwen/6-07082020110802.html")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetPost(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}

func TestRfa(t *testing.T) {
	p := PostFactory("https://www.rfa.org/mandarin/Xinwen/6-07082020110802.html")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := Rfa(p)
	fmt.Println(tc)
}
