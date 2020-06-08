package fetcher

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func FetchBoxun(url string) (*fetcher.ThePost, error) {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}
	domain := "www.boxun.com"
	site := "www.boxun.com"
	title := ThisGetTitle(rawBody)
	// get contents
	body, err := FmtBodyBoxun(rawBody)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}
	a := filepath.Base(url)

	date := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", a[:4], a[4:6], a[6:8], a[8:10], a[10:12], "00")

	post := fetcher.ThePost{Site: site, Domain: domain, URL: url, Title: title, Body: body, Date: date}

	return &post, nil
}

func ThisGetTitle(raw string) string {
	var a = regexp.MustCompile(`(?m)<title>(.*?)</title>`)
	rt := a.FindStringSubmatch(raw)
	if rt != nil {
		s := rt[1]
		return gears.ConvertToUtf8(s, "gbk", "utf-8")

	} else {
		return ""

	}
}

func FetchBoxunUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<li><a href="(/news/\w*/.*?)" target=_blank>`)
	lst := reLink.FindAllStringSubmatch(rawBody, -1)
	if lst == nil {
		fmt.Printf("\n[-] fetcher.FetchBoxunUrls(%s) regex matched nothing.\n", url)
		return nil
	} else {
		for _, v := range reLink.FindAllStringSubmatch(rawBody, -1) {
			ret_lst = append(ret_lst, v[1])
		}
		ret_lst = gears.StrSliceDeDupl(ret_lst)
	}

	return ret_lst
}

// FmtBodyBoxun focus on dwnews, it can extract raw body string via regexp and then, format the news body to markdowned string.
func FmtBodyBoxun(rawBody string) (string, error) {
	if rawBody == "" {
		return "", errors.New("[-] FmtBodyBoxun() parameter is nil!")
	}
	var ps []string
	var body string
	var re = regexp.MustCompile(`(?m)<!--bodystart-->([^\^]*)<!--bodyend-->`)
	for _, v := range re.FindAllStringSubmatch(rawBody, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return "", errors.New("[-] fetcher.FmtBodyBoxun() Error: regex matched nothing.")
	} else {
		for _, p := range ps {
			body += p + "  \n"
		}
	}
	a := regexp.MustCompile(`<BR>`)
	bodySlice := a.Split(body, -1)
	for _, v := range bodySlice {
		re = regexp.MustCompile(`&nbsp;`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<div(.*?)</div>`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<img(.*?)>`)
		v = re.ReplaceAllString(v, "")
		body += v + "  \n"
	}
	return body, nil
}
