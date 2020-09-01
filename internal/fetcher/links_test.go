package fetcher

import (
	"net/url"
	"testing"
)

func TestKickOutLinksMatchPath(t *testing.T) {
	// link := "https://www.dwnews.com/%E8%A7%86%E8%A7%89/60202427/"
	beKick := []string{"https://www.dwnews.com/%E8%A7%86%E8%A7%89/60202427/"}
	// path := "/" + url.QueryEscape("视觉") + "/"
	path := url.QueryEscape("视觉")
	KickOutLinksMatchPath(&beKick, path)
	if len(beKick) != 0 {
		t.Errorf("want: len(beKick) == 0, got: len(beKick) == %d", len(beKick))
	}
}

func TestSetLinks(t *testing.T) {
	u, err := url.Parse("https://news.ltn.com.tw/list/breakingnews")
	if err != nil {
		t.Errorf("Url Parse fail!\n%s", err)
	}
	var f = &Fetcher{
		Entrance: u,
	}
	f.SetLinks()
	assertLinks := []string{
		"https://news.ltn.com.tw/news/society/breakingnews/3278253",
		"https://news.ltn.com.tw/news/society/breakingnews/3278250",
		"https://news.ltn.com.tw/news/politics/breakingnews/3278225",
		"https://news.ltn.com.tw/news/politics/breakingnews/3278170",
	}
	shot := 0
	for _, link := range f.Links {
		for _, v := range assertLinks {
			if link == v {
				shot++
			}
		}
	}
	if shot == 0 {
		t.Errorf("want: %v, got: %v", len(assertLinks), shot)
	}
}
