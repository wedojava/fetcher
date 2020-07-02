package fetcher

import (
	"testing"
	"time"
)

func TestSetTitle(t *testing.T) {
	p := PostFactory("https://www.dwnews.com/%E4%B8%AD%E5%9B%BD/60202347")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc error: %v", err)
	}
	p.DOC, p.Raw = doc, raw
	err = p.SetTitle()
	if err != nil {
		t.Errorf("SetTitle err: %v", err)
	}
	want := "【港版国安法】条文解读：再有大动乱 武警或可跨过深圳湾"
	if p.Title != want {
		t.Errorf("want: %s, got: %s", p.Title, want)
	}
}
