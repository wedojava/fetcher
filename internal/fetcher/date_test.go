package fetcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/wedojava/fetcher/internal/htmldoc"
)

func TestVoaDateInNode(t *testing.T) {
	p := PostFactory("https://www.dwnews.com/%E9%A6%99%E6%B8%AF/60201980/")
	_, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetDOC error: %v", err)
	}
	p.DOC = doc
	if err := p.VoaDateInNode(); err != nil {
		t.Errorf("VoaDateInNode error: %v", err)
	}
	fmt.Println(p.Date)
}

func TestSetDate(t *testing.T) {
	f := FetcherFactory("https://www.dwnews.com")
	if err := f.SetLinks(); err != nil {
		t.Errorf("SetLinks error: %v", err)
	}
	for _, link := range f.Links {
		p := PostFactory(link)
		_, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
		if err != nil {
			t.Errorf("GetDOC error: %v", err)
		}
		p.DOC = doc
		if err := p.VoaDateInNode(); err != nil {
			t.Errorf("VoaDateInNode error: %v", err)
		}
		fmt.Println(p.Date)

	}
}

func TestBoxunDateInUrl(t *testing.T) {
	p := PostFactory("https://www.boxun.com/news/gb/intl/2020/06/202006302339.shtml")
	if err := p.BoxunDateInUrl(); err != nil {
		t.Errorf("BoxunDateInUrl() invoked error: %v", err)
	}
	want := "2020-06-30T23:39:00Z"
	if p.Date != want {
		t.Errorf("want: %v, got: %v", want, p.Date)
	}
	fmt.Println(p.Date)
}
