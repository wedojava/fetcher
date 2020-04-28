package fetcher

import (
	"fmt"
	"testing"
)

func checkFetch(t *testing.T, _got, _want string) {
	t.Helper()
	if _got != _want {
		t.Errorf("\ngot %v\nwant %v\n", _got, _want)
	}
}

func TestFetch(t *testing.T) {
	t.Run("test get title and body: ", func(t *testing.T) {
		got, _ := Fetch("https://www.dwnews.com/%E5%85%A8%E7%90%83/60176216/%E6%96%B0%E5%86%A0%E8%82%BA%E7%82%8E%E6%9C%80%E6%96%B0%E7%96%AB%E6%83%85%E5%85%A8%E7%90%83%E7%A1%AE%E8%AF%8A%E9%80%BE256%E4%B8%87%E4%BE%8B%E7%BE%8E%E5%9B%BD%E7%A1%AE%E8%AF%8A82%E4%B8%87%E4%BE%8B")
		wantTitle := "【新冠肺炎·最新疫情】全球确诊逾256万例 美国确诊82万例"
		wantSite := "@dwnewsofficial"
		wantDate := "2020-04-22T08:55:02+08:00"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Site, wantSite)
		checkFetch(t, got.Date, wantDate)
	})
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchUrls("https://www.dwnews.com")
		want := 43
		if len(got) == want {
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}
