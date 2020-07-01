package fetcher

import (
	"errors"
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
	}
	return nil
}

func (p *Post) FmtBody(f func(body []byte) (string, error)) error {
	if p.DOC == nil {
		return errors.New("[-] there is no DOC object to get and format.")
	}
	b, err := f(p.Raw)
	if err != nil {
		return err
	}
	p.Body = "#" + p.Filename[:strings.LastIndex(p.Filename, ".")] + "\n\n" + b
	return nil
}

func Boxun(body []byte) (string, error) {
	var ps []string
	var _body string
	var re = regexp.MustCompile(`(?m)&nbsp;&nbsp;&nbsp;&nbsp;(.*?)<BR>`)
	vs := re.FindAllSubmatch(body, -1)
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
