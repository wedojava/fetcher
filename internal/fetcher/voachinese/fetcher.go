package fetcher

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func FetchVoa(url string) (*fetcher.ThePost, error) {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
	domain := gears.HttpGetDomain(url)
	site := gears.HttpGetSiteViaTwitterJS(rawBody)
	title := gears.HttpGetTitleViaTwitterJS(rawBody)
	// get contents
	body, err := FmtBodyVoa(rawBody)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
	date := gears.HttpGetDateByHeader(rawBody)

	post := fetcher.ThePost{Site: site, Domain: domain, URL: url, Title: title, Body: body, Date: date}

	return &post, nil
}

func FetchVoaUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<a\s+href\s*=\s*"(?P<links>/a/.*?)"\s*>`)
	lst := reLink.FindAllStringSubmatch(rawBody, -1)
	if lst == nil {
		fmt.Println("[-] fetcher.FetchVoaUrls() regex matched nothing.")
		return nil
	} else {
		for _, v := range reLink.FindAllStringSubmatch(rawBody, -1) {
			ret_lst = append(ret_lst, v[1])
		}
		ret_lst = gears.StrSliceDeDupl(ret_lst)
	}

	return ret_lst
}

// FmtBodyRfa focus on dwnews, it can extract raw body string via regexp and then, format the news body to markdowned string.
func FmtBodyVoa(rawBody string) (string, error) {
	var ps []string
	var body string
	var reContent = regexp.MustCompile(`(?m)<p>(?P<content>.*?)</p>`)
	for _, v := range reContent.FindAllStringSubmatch(rawBody, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return "", errors.New("[-] fetcher.FmtBodyVoa() Error: regex matched nothing.")
	} else {
		for _, p := range ps {
			body += p + "  \n"
		}
	}

	return body, nil
}
