package fetcher

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func FetchRfa(url string) (*fetcher.ThePost, error) {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}
	domain := gears.HttpGetDomain(url)
	site := gears.HttpGetSiteViaTwitterJS(rawBody)
	title := gears.HttpGetTitleViaTwitterJS(rawBody)
	// get contents
	body, err := FmtBodyRfa(rawBody)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}
	date := gears.HttpGetDateByHeader(rawBody)

	post := fetcher.ThePost{Site: site, Domain: domain, URL: url, Title: title, Body: body, Date: date}

	return &post, nil
}

func FetchRfaUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<a href\s*=\s*"\s*(https://www.rfa.org/.*?.html)\s*"\s*>`)
	lst := reLink.FindAllStringSubmatch(rawBody, -1)
	if lst == nil {
		fmt.Printf("\n[-] fetcher.FetchRfaUrls(%s) regex matched nothing.\n", url)
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
func FmtBodyRfa(rawBody string) (string, error) {
	if rawBody == "" {
		return "", errors.New("[-] FmtBodyVoa() parameter is nil!")
	}
	var ps []string
	var body string
	var reContent = regexp.MustCompile(`(?m)<p.*?>(?P<content>.*?)</p>`)
	for _, v := range reContent.FindAllStringSubmatch(rawBody, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return "", errors.New("[-] fetcher.FmtBodyRfa() Error: regex matched nothing.")
	} else {
		for _, p := range ps {
			body += p + "  \n"
		}
	}

	return body, nil
}
