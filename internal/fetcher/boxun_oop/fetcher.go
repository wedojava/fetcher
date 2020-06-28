package fetcher

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

type PostBoxun struct {
	fetcher.ThePost
}

func (p *PostBoxun) FetchBoxun() error {
	// set contents
	doc, err := fetcher.GetDOC(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.DOC = doc
	raw, err := fetcher.GetRaw(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.Raw = raw
	// set Date
	url, err := url.Parse(p.URL)
	if err != nil {
		return err
	}
	a := filepath.Base(url.Path)

	p.Date = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", a[:4], a[4:6], a[6:8], a[8:10], a[10:12], "00")
	// set Title
	err = p.GetTitle()
	if err != nil {
		return err
	}
	p.Title = gears.ConvertToUtf8(p.Title, "gbk", "utf8")
	p.Title = strings.TrimSpace(p.Title)
	// fmt body
	err = p.FmtBodyBoxun()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// TODO: use links to implement this func
func FetchBoxunUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<a\shref\s?=\s?"(?P<href>/.{2}/\d{8}/.+?)".*?>`)
	lst := reLink.FindAllStringSubmatch(rawBody, -1)
	if lst == nil {
		fmt.Printf("\n[-] fetcher.FetchBoxunUrls(%s) regex matched nothing.\n", url)
		return nil
	} else {
		for _, v := range reLink.FindAllStringSubmatch(rawBody, -1) {
			ret_lst = append(ret_lst, "https://www.dwnews.com"+v[1])
		}
		ret_lst = gears.StrSliceDeDupl(ret_lst)
	}

	return ret_lst
}

// FmtBodyBoxun can extract html body and fmt this string by doc of body.
func (p *PostBoxun) FmtDocBoxun() error {
	if p.DOC == nil {
		return errors.New("[-] FmtBodyBoxun() parameter is nil!")
	}
	var ps []string
	var body string
	var re = regexp.MustCompile(`(?m)<!--bodystart-->([^\^]*)<!--bodyend-->`)
	for _, v := range re.FindAllStringSubmatch(p.Body, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return errors.New("[-] fetcher.FmtBodyRfa() Error: regex matched nothing.")
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
	p.Body = body
	return nil
}

// FmtBodyBoxun can extract html body and fmt this string by raw of body.
func (p *PostBoxun) FmtBodyBoxun() error {
	if p.Raw == nil {
		return errors.New("[-] FmtBodyBoxun() parameter is nil!")
	}
	var ps []string
	var body string
	var re = regexp.MustCompile(`(?m)<!--bodystart-->([^\^]*)<!--bodyend-->`)
	for _, v := range re.FindAllStringSubmatch(p.Body, -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return errors.New("[-] fetcher.FmtBodyRfa() Error: regex matched nothing.")
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
	p.Body = body
	return nil
}
