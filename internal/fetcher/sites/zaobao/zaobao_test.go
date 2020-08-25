package zaobao

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

func TestSetDate(t *testing.T) {
	p := PostFactory("https://www.zaobao.com/realtime/world/story20200825-1079575")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetDate(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	want := "2020-08-25T09:42:32+08:00"
	if p.Date != want {
		t.Errorf("\ngot: %v\nwant: %v", p.Date, want)
	}
}

func TestSetTitle(t *testing.T) {
	p := PostFactory("https://www.zaobao.com/realtime/world/story20200825-1079575")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := SetTitle(p); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	want := "韩首都圈学校全面线上上课两周"
	if p.Title != want {
		t.Errorf("\ngot: %v\nwant: %v", p.Title, want)
	}
}

func TestZaobao(t *testing.T) {
	p := PostFactory("https://www.zaobao.com/realtime/world/story20200825-1079575")
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	tc, err := Zaobao(p)
	fmt.Println(tc)
}

func TestSetPost(t *testing.T) {
	p := PostFactory("https://www.zaobao.com/realtime/world/story20200825-1079575")
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
