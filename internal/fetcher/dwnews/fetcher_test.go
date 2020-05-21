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
		got, _ := FetchDwnews("https://www.dwnews.com/%E5%85%A8%E7%90%83/60176216/%E6%96%B0%E5%86%A0%E8%82%BA%E7%82%8E%E6%9C%80%E6%96%B0%E7%96%AB%E6%83%85%E5%85%A8%E7%90%83%E7%A1%AE%E8%AF%8A%E9%80%BE256%E4%B8%87%E4%BE%8B%E7%BE%8E%E5%9B%BD%E7%A1%AE%E8%AF%8A82%E4%B8%87%E4%BE%8B")
		wantTitle := "【新冠肺炎·最新疫情】全球确诊逾256万例 美国确诊82万例"
		wantDomain := "www.dwnews.com"
		wantSite := "@dwnewsofficial"
		wantDate := "2020-04-22T08:55:02+08:00"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Domain, wantDomain)
		checkFetch(t, got.Site, wantSite)
		checkFetch(t, got.Date, wantDate)
	})
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchDwnewsUrls("https://www.dwnews.com")
		want := 42
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
		summary := "台湾总统蔡英文正式展开第二任期，美国国务院5月20日表示，批准以1.8亿美元向台湾出售18枚重量级鱼雷。此举料进一步加剧华盛顿与北京已经紧张的关系。"
		if strings.Contains(body, summary) {
			fmt.Println("Success")
		} else {
			fmt.Println("Cannot fetch summary correctly.")
		}

	})
}
