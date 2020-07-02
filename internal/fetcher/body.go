package fetcher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/wedojava/gears"
)

func (p *Post) SetBody() error {
	switch p.Domain {
	case "www.boxun.com":
		if err := p.FmtBody(Boxun); err != nil {
			return err
		}
	case "www.dwnews.com":
		if err := p.FmtBody(Dwnews); err != nil {
			return err
		}
	case "www.voachinese.com":
		if err := p.FmtBody(Voa); err != nil {
			return err
		}
	}
	return nil
}

func (p *Post) FmtBody(f func(post *Post) (string, error)) error {
	if p.DOC == nil {
		return errors.New("[-] there is no DOC object to get and format.")
	}
	b, err := f(p)
	if err != nil {
		return err
	}
	p.Body = "#" + p.Filename[:strings.LastIndex(p.Filename, ".")] + "\n\n" + b
	return nil
}

func Boxun(p *Post) (string, error) {
	var ps []string
	var _body string
	var re = regexp.MustCompile(`(?m)<!--bodystart-->([^\^]*)<!--bodyend-->`)
	vs := re.FindAllSubmatch(p.Raw, -1)
	if vs == nil {
		var re = regexp.MustCompile(`(?m)&nbsp;&nbsp;&nbsp;&nbsp;(.*?)<BR>`)
		vs = re.FindAllSubmatch(p.Raw, -1)
	}
	for _, v := range vs {
		ps = append(ps, string(v[1]))
	}
	if len(ps) == 0 {
		return "", errors.New("[-] Boxun() match nothing from body.")
	} else {
		for _, p := range ps {
			_body += p + "  \n"
		}
	}
	a := regexp.MustCompile(`<BR>`)
	bodySlice := a.Split(string(_body), -1)
	_body = ""
	for _, v := range bodySlice {
		re = regexp.MustCompile(`<table([^\^]*)</table>`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`&nbsp;`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<div(.*?)</div>`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<img(.*?)>`)
		v = re.ReplaceAllString(v, "")
		re = regexp.MustCompile(`<P>`)
		v = re.ReplaceAllString(v, "")
		if strings.TrimSpace(v) == "" {
			continue
		}
		_body += v + "  \n"
	}
	if err := gears.ConvertToUtf8(&_body, "gbk", "utf-8"); err != nil {
		return "", err
	}
	return _body, nil
}

func Dwnews(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	nodes := ElementsByTagName(doc, "article")
	if len(nodes) == 0 {
		return "", errors.New("[-] Got nothing from: " + p.URL.String())
	}
	articleDoc := nodes[0].FirstChild
	plist := ElementsByTagName(articleDoc, "p")
	if articleDoc.FirstChild.Data == "div" {
		body += fmt.Sprintf("\n > %s  \n", plist[0].FirstChild.Data)
	}
	for _, v := range plist { // the last item is `推荐阅读：`
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.FirstChild != nil && v.FirstChild.Data == "strong" {
			body += fmt.Sprintf("** %s **  \n", v.FirstChild.FirstChild.Data)
		} else {
			body += v.FirstChild.Data + "  \n"
		}
		body += v.FirstChild.Data + "  \n"
	}
	body = strings.ReplaceAll(body, "strong", "")
	body = strings.ReplaceAll(body, "推荐阅读：", "")
	return body, nil
}

func Voa(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	articleDoc := ElementsByTagAndClass(doc, "div", "wsw")
	plist := ElementsByTagName(articleDoc[0], "p")
	for _, v := range plist {
		body += v.FirstChild.Data + "  \n"
	}
	return body, nil
}
