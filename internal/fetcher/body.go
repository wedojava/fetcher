package fetcher

import (
	"errors"
	"regexp"
	"strings"
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
	p.Body = b
	return nil
}

func Boxun(body []byte) (string, error) {
	var ps []string
	var _body string
	var re = regexp.MustCompile(`(?m)<!--bodystart-->([^\^]*)<!--bodyend-->`)
	// TODO: replace implement from string dealing to []byte dealing
	for _, v := range re.FindAllStringSubmatch(string(body), -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		return "", errors.New("[-] Boxun() regex error: want match body but matched nothing.")
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
		if strings.TrimSpace(v) == "" {
			continue
		}
		_body += v + "  \n"
	}
	return _body, nil
}
