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
		got, _ := FetchBoxun("https://boxun.com/news/gb/intl/2020/06/202006081127.shtml")
		wantTitle := "无公正 不停歇：欧洲各大城市加入反警暴浪潮"
		wantDomain := "www.boxun.com"
		wantSite := "www.boxun.com"
		wantDate := "2020-06-08T11:27:00Z"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Domain, wantDomain)
		checkFetch(t, got.Site, wantSite)
		checkFetch(t, got.Date, wantDate)
	})
}

func TestFmtBodyBoxun(t *testing.T) {
	rawBody, err := gears.HttpGetBody("https://boxun.com/news/gb/pubvp/2020/06/202006110029.shtml", 10)
	// rawBody, err := gears.HttpGetBody("https://boxun.com/news/gb/intl/2020/06/202006081127.shtml", 10)
	if err != nil {
		t.Fatal(err)
	}
	got, err := FmtBodyBoxun(rawBody)
	if strings.Contains(got, "還有「犯我強漢者，雖遠必誅」，「壯志飢餐胡虜肉」等，都是漢人愛國將領的名言。西方國家要記住這些話。台灣人被選上總統，對外省權貴來說，也是「犯漢」，一定要用各種手段進行報復。") {
		// if strings.Contains(got, "意大利罗马的人民广场也以8分钟的寂静来悼念弗洛伊德，选择8分钟，是因为弗洛伊德被白人警察用膝盖压住脖子时间超过8分钟。在集体默哀之后，人们继续呼喊反种族歧视的口号。受到非法难民浪潮等社会政治经济问题的冲击，意大利向民粹偏移，本次罗马游行的人群当中不少都是移民或者非法入境难民的身份，他们表示，在意大利正常生活，真的很难。") {
		fmt.Print("Test pass.")
	}
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchBoxunUrls("https://www.boxun.com/rolling.shtml")
		want := 179
		if len(got) == want {
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}
