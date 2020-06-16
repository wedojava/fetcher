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
	rawBody, err := gears.HttpGetBody("https://www.rfa.org/mandarin/Xinwen/1-06142020094439.html", 10)
	// rawBody, err := gears.HttpGetBody("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/wy-05052020131409.html", 10)
	if err != nil {
		t.Fatal(err)
	}
	got, err := FmtBodyRfa(rawBody)
	if strings.Contains(got, "韩国总统府青瓦台周日凌晨紧急召开国安会议，外交部长官康京和、统一部长官金炼铁、国防部长官郑景斗、国家情报院院长徐薰等人就当前半岛局势及解决方案进行磋商。韩国统一部当天也呼吁双方努力遵守所有协议，对“当前情况十分严重”的半岛局势表达担忧。韩国国防部当天则表示，正保持高度戒备，密切关注朝军动向") {
		fmt.Println("done.")
	} else {
		t.Errorf("got: %v", got)
	}

}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchRfaUrls("https://www.rfa.org/mandarin/")
		want := 34
		if len(got) == want {
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}
