package fetcher

import (
	"testing"
	"time"
)

func TestSetTitle(t *testing.T) {
	// p := PostFactory("https://www.dwnews.com/%E4%B8%AD%E5%9B%BD/60202347")
	p := PostFactory("https://www.boxun.com/news/gb/intl/2020/07/202007041307.shtml")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc error: %v", err)
	}
	p.DOC, p.Raw = doc, raw
	if err = p.SetTitle(); err != nil {
		t.Errorf("SetTitle err: %v", err)
	}
	want := "朱万利：郭文贵起诉案件进展，美国对其金融诈骗立案调查"
	if p.Title != want {
		t.Errorf("got: %s, want: %s", p.Title, want)
	}
}
