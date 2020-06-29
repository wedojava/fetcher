package fetcher

import (
	"fmt"
	"testing"
)

var testPost PostBoxun

// ThePost.Site: "www.boxun.com",
// Domain:       "www.boxun.com",
// URL:          "https://boxun.com/news/gb/intl/2020/06/202006081127.shtml",
// Title:        "",
// Body:         "",
// Date:         "",

func checkFetch(t *testing.T, _got, _want string) {
	t.Helper()
	if _got != _want {
		t.Errorf("\ngot %v\nwant %v\n", _got, _want)
	}
}

func TestFetchTitleAndDate(t *testing.T) {
	testPost.Domain = "www.boxun.com"
	testPost.Entrance = "www.boxun.com"
	testPost.URL = "https://boxun.com/news/gb/intl/2020/06/202006081127.shtml"
	testPost.FetchBoxun()
	t.Run("test get title and body: ", func(t *testing.T) {
		wantTitle := "无公正 不停歇：欧洲各大城市加入反警暴浪潮"
		wantDate := "2020-06-08T11:27:00Z"
		checkFetch(t, testPost.Title, wantTitle)
		checkFetch(t, testPost.Date, wantDate)
	})
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchBoxunUrls("https://news.boxun.com/")
		want := 43
		if len(got) == want {
			fmt.Println(got[0])
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}
