package fetcher

import (
	"fmt"
	"testing"
	"time"
)

func TestVoaDateInNode(t *testing.T) {
	p := PostFactory("https://www.dwnews.com/%E9%A6%99%E6%B8%AF/60201980/")
	_, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
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
		_, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
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
