package fetcher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wedojava/gears"
)

func checkFetch(t *testing.T, _got, _want string) {
	t.Helper()
	if _got != _want {
		t.Errorf("\ngot %v\nwant %v\n", _got, _want)
	}
}

func TestFetch(t *testing.T) {
	t.Run("test get title and body: ", func(t *testing.T) {
		// got, _ := FetchDwnews("https://www.dwnews.com/全球/60181030/")
		got, _ := FetchDwnews("https://www.dwnews.com/中国/60179204")
		wantTitle := "【对台军售】美批准对台售18枚MK48重型鱼雷"
		wantDomain := "www.dwnews.com"
		wantSite := "@dwnewsofficial"
		wantDate := "2020-05-21T09:31:02+08:00"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Domain, wantDomain)
		checkFetch(t, got.Site, wantSite)
		checkFetch(t, got.Date, wantDate)
	})
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchDwnewsUrls("https://www.dwnews.com")
		want := 43
		if len(got) == want {
			fmt.Println(got[0])
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}

func TestFmtBodyDwnews(t *testing.T) {
	raw, _ := gears.HttpGetBody("https://www.dwnews.com/中国/60179204", 10)
	body, _ := FmtBodyDwnews(raw)
	t.Run("test fetch summary then fmt it: ", func(t *testing.T) {
		summary := "> 台湾总统蔡英文正式展开第二任期，美国国务院5月20日表示，批准以1.8亿美元向台湾出售18枚重量级鱼雷。此举料进一步加剧华盛顿与北京已经紧张的关系。"
		if strings.Contains(body, summary) {
			fmt.Println("Success")
		} else {
			fmt.Println("Cannot fetch summary correctly.")
		}

	})
}
