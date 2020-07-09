package fetcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/wedojava/fetcher/internal/htmldoc"
)

func TestSetAndSavePost(t *testing.T) {
	// p := PostFactory("https://www.dwnews.com/经济/60203253")
	p := PostFactory("https://www.dwnews.com/经济/60203034") // The wrong one
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := p.SetPost(); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	if err := p.SavePost(); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}

func TestTreatPost(t *testing.T) {
	tcs := []string{
		"https://www.boxun.com/news/gb/taiwan/2020/07/202007091815.shtml",
		"https://www.dwnews.com/经济/60203253",
		"https://www.dwnews.com/全球/60203234",
		"https://www.voachinese.com/a/S-Korea-Says-US-Sees-Importance-Of-N-Korea-Talks-Despite-Tension-20200709/5496028.html",
	}
	for _, tc := range tcs {
		p := PostFactory(tc)
		p.TreatPost()
	}
}
