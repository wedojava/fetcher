package fetcher

import (
	"fmt"
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
	u, err := url.Parse("https://www.zaobao.com/realtime/world")
	if err != nil {
		t.Errorf("Url Parse fail!\n%s", err)
	}
	var f = &Fetcher{
		Entrance: u,
	}
	f.SetLinks()
	assertLink := "https://www.zaobao.com/realtime/world/story20200825-1079575"
	shot := 0
	for _, link := range f.Links {
		fmt.Println(link)
		if link == assertLink {
			shot++
		}
	}
	if shot == 0 {
		t.Errorf("want: %v, got: %v", 1, shot)
	}
}
