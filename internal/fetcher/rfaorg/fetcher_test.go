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
		got, _ := FetchRfa("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/wy-05052020131409.html")
		wantTitle := "《病毒往事》爆红网络 中国宣传、智囊机构加入舆论战"
		wantDomain := "www.rfa.org"
		wantSite := "@RadioFreeAsia"
		wantDate := "2020-05-05T19:44:28Z"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Domain, wantDomain)
		checkFetch(t, got.Site, wantSite)
		checkFetch(t, got.Date, wantDate)
	})
}

func TestFmtBodyRfa(t *testing.T) {
	rawBody, err := gears.HttpGetBody("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/wy-05052020131409.html", 10)
	if err != nil {
		t.Fatal(err)
	}
	got, err := FmtBodyRfa(rawBody)
	if strings.Contains(got, "视频只有1分45秒长，画面简单，是一个代表美国的“自由女神像”和一个代表中国的带着口罩的兵马俑形象在作英文对话。主要内容是，从去年12月到今年4月，中国及早发现并报告了病毒，而美方防疫不力，造成新冠病毒在美国大流行") {
		fmt.Print("Test pass.")
	}
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
