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
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}
	domain := gears.HttpGetDomain(url)
	site := gears.HttpGetSiteViaTwitterJS(rawBody)
	title := ThisGetTitle(rawBody)
	// get contents
	body, err := FmtBodyVoa(rawBody)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}
	date := ThisGetDate(rawBody)
	// date = date[:10] + "T" + date[11:]

	post := fetcher.ThePost{Site: site, Domain: domain, URL: url, Title: title, Body: body, Date: date}

	return &post, nil
}

func ThisGetDate(rawBody string) string {
	if rawBody == "" {
		return ""
	}
	var a = regexp.MustCompile(`(?m)<time\s+datetime="(?P<date>.*?)">\n*.*\n*</time>`)
	rt := a.FindStringSubmatch(rawBody)
	if rt != nil {
		return rt[1]

	} else {
		return ""

	}
}

func ThisGetTitle(rawBody string) string {
	var a = regexp.MustCompile(`(?m)<title>(?P<title>.*?)</title>`)
	rt := a.FindStringSubmatch(rawBody)
	if rt != nil {
		return rt[1]

	} else {
		return ""

	}
}

func FetchVoaUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<a\s+href\s*=\s*"(?P<links>/a/.*-.*.html)"\s*>`)
	lst := reLink.FindAllStringSubmatch(rawBody, -1)
	if lst == nil {
		fmt.Printf("\n[-] fetcher.FetchVoaUrls(%s) regex matched nothing.\n", url)
		return nil
	} else {
		for _, v := range reLink.FindAllStringSubmatch(rawBody, -1) {
			ret_lst = append(ret_lst, "https://www.voachinese.com"+v[1])
		}
		ret_lst = gears.StrSliceDeDupl(ret_lst)
	}

	return ret_lst
}

// FmtBodyVoa focus on voa, it can extract raw body string via regexp and then, format the news body to markdowned string.
func FmtBodyVoa(rawBody string) (string, error) {
	if rawBody == "" {
		return "", errors.New("\n[-] FmtBodyVoa() parameter is nil!\n")
	}
	var ps []string
	var body string
	var reContent = regexp.MustCompile(`(?m)<p>(?P<content>.*?)</p>`)
	for _, v := range reContent.FindAllStringSubmatch(rawBody, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		if regexp.MustCompile(`(?m)<video.*?>`).FindAllString(rawBody, -1) != nil {
			return "", errors.New("\n[-] fetcher.FmtBodyVoa() Error: this is a video page.\n")
		}
		return "", errors.New("\n[-] fetcher.FmtBodyVoa() Error: regex matched nothing.\n")
	} else {
		for _, p := range ps {
			body += p + "  \n"
		}
	}

	return body, nil
}
