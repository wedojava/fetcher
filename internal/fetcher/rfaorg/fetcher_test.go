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
		got, _ := FetchRfa("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/wy-05052020131409.html")
		wantTitle := "《病毒往事》爆红网络 中国宣传、智囊机构加入舆论战"
		wantDomain := "www.rfa.org"
		wantSite := "@RadioFreeAsia"
		// wantDate := "2020-04-22T08:55:02+08:00"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Domain, wantDomain)
		checkFetch(t, got.Site, wantSite)
		// checkFetch(t, got.Date, wantDate)
	})
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchRfaUrls("https://www.dwnews.com")
		want := 43
		if len(got) == want {
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}
