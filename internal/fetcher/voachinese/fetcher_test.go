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
		got, _ := FetchVoa("https://www.voachinese.com/a/pompeo-warns-china-over-interference-with-us-journalists-in-hk/5423711.html")
		wantTitle := "蓬佩奥警告中国不要干涉在香港的美国记者的工作"
		wantDomain := "www.voachinese.com"
		// wantSite := "@RadioFreeAsia"
		wantDate := "2020-05-17 22:44:15Z"
		checkFetch(t, got.Title, wantTitle)
		checkFetch(t, got.Domain, wantDomain)
		// checkFetch(t, got.Site, wantSite)
		checkFetch(t, got.Date, wantDate)
	})
}

func TestFmtBodyRfa(t *testing.T) {
	rawBody, err := gears.HttpGetBody("https://www.voachinese.com/a/pompeo-warns-china-over-interference-with-us-journalists-in-hk/5423711.html", 10)
	if err != nil {
		t.Fatal(err)
	}
	got, err := FmtBodyVoa(rawBody)
	if strings.Contains(got, "蓬佩奥在一份声明中说：“这些记者是自由媒体的成员，而不是宣传干部，他们的有价值的报道可以告知中国公民和世界。”") {
		fmt.Print("Test pass.")
	}
}

func TestFetchUrls(t *testing.T) {
	t.Run("get urls count.", func(t *testing.T) {
		got := FetchVoaUrls("https://www.voachinese.com")
		want := 89
		if len(got) == want {
			fmt.Print("Test pass.")
		} else {
			t.Errorf("\nGot %v\n Want %v", len(got), want)
		}

	})
}
