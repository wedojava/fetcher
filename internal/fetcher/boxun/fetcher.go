package fetcher

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
	title := TitleBoxun(rawBody)
	// get contents
	body, err := FmtBodyBoxun(rawBody)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return nil, err
	}

	date := DateBoxun(url)

	post := fetcher.ThePost{Site: site, Domain: domain, URL: url, Title: title, Body: body, Date: date}

	return &post, nil
}

func DateBoxun(url string) string {
	rawdate := filepath.Base(url)
	var Y, M, D, hh, mm int
	var err error
	if Y, err = strconv.Atoi(rawdate[:4]); err != nil {
		Y = 2000
		fmt.Println(rawdate[:4], "is not a Year., set it to 2000")
	}
	if M, err = strconv.Atoi(rawdate[4:6]); err != nil {
		M = 0
		fmt.Println(rawdate[4:6], "is not a Month., set it to 0")
	}
	if D, err = strconv.Atoi(rawdate[6:8]); err != nil {
		D = 0
		fmt.Println(rawdate[6:8], "is not a Date., set it to 0")
	}
	if hh, err = strconv.Atoi(rawdate[8:10]); err != nil {
		hh = 0
		fmt.Println(rawdate[8:10], "is not a Hour., set it to 0")
	}
	if mm, err = strconv.Atoi(rawdate[10:12]); err != nil {
		mm = 0
		fmt.Println(rawdate[10:12], "is not a Minute., set it to 0")
	}
	if err != nil {
		fmt.Println("err date fetch from url: ", url)
	}
	if Y < 1999 || Y > 2499 {
		Y = 2000
		// fmt.Println(Y, "is not a integer of Year, set it to 2000")
	}
	if M < 0 || M >= 12 {
		M = 12
		// fmt.Println(M, "is not a integer of Month, set it to 12")
	}
	if D < 0 || D >= 31 {
		D = 31
		// fmt.Println(D, "is not a integer of Date, set it to 31")
	}
	if hh < 0 || hh > 23 {
		hh = 23
		// fmt.Println(hh, "is not a integer of Hour, set it to 23")
	}
	if mm < 0 || mm > 59 {
		mm = 59
		// fmt.Println("err date fetch from url: ", url)
		// fmt.Println(mm, "is not a integer of Minute, set it to 59")
	}
	return fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02dZ", Y, M, D, hh, mm, 0)

}

func TitleBoxun(raw string) string {
	var a = regexp.MustCompile(`(?m)<title>(.*?)</title>`)
	rt := a.FindStringSubmatch(raw)
	if rt != nil {
		s := strings.TrimSpace(rt[1])
		s = gears.ConvertToUtf8(s, "gbk", "utf-8")
		fetcher.ReplaceIllegalChar(&s)
		return s
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
			ret_lst = append(ret_lst, "https://boxun.com"+v[1])
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
	body = ""
	for _, v := range bodySlice {
		re = regexp.MustCompile(`&nbsp;`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<div(.*?)</div>`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<img(.*?)>`)
		v = re.ReplaceAllString(v, "")
		if strings.TrimSpace(v) == "" {
			continue
		}
		body += v + "  \n"
	}
	body = gears.ConvertToUtf8(body, "gbk", "utf-8")
	return body, nil
}
