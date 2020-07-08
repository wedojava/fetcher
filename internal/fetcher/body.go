package fetcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wedojava/gears"
	"golang.org/x/net/html"
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
	case "www.rfa.org":
		if err := p.FmtBody(Rfa); err != nil {
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
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := ElementsByTagAndClass(doc, "td", "F11")
	if len(nodes) == 0 {
		return "", errors.New("[-] There is no tag named `<td class=F11>` from: " + p.URL.String())
	}
	blist := ElementsNextByTag(nodes[0], "br")
	for _, b := range blist {
		if b.Type != html.TextNode || b.Data == "" {
			continue
		} else {
			body += strings.ReplaceAll(b.Data, "\u00a0", "")
		}
	}
	gears.ConvertToUtf8(&body, "gbk", "utf8")
	return body, nil
}

func Dwnews(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := ElementsByTagName(doc, "article")
	if len(nodes) == 0 {
		return "", errors.New("[-] There is no tag named `<article>` from: " + p.URL.String())
	}
	articleDoc := nodes[0].FirstChild
	plist := ElementsByTagName(articleDoc, "p")
	if articleDoc.FirstChild.Data == "div" { // to fetch the summary block
		body += fmt.Sprintf("\n > %s  \n", plist[0].FirstChild.Data)
	}
	for _, v := range plist { // the last item is `推荐阅读：`
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.FirstChild != nil && v.FirstChild.Data == "strong" {
			body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.FirstChild.Data)
		} else {
			body += v.FirstChild.Data + "  \n"
		}
	}
	body = strings.ReplaceAll(body, "strong", "")
	body = strings.ReplaceAll(body, "** 推荐阅读： **", "")
	return body, nil
}

func Voa(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := ElementsByTagAndClass(doc, "div", "wsw")
	if len(nodes) == 0 {
		return "", errors.New(`[-] There is no element match '<div class="wsw">'`)
	}
	plist := ElementsByTagName(nodes[0], "p")
	for _, v := range plist {
		body += v.FirstChild.Data + "  \n"
	}
	body = strings.ReplaceAll(body, "strong  \n", "")
	body = strings.ReplaceAll(body, "span  \n", "")
	body = strings.ReplaceAll(body, "br  \n", "")
	return body, nil
}

func Rfa(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := ElementsByTagAndId(doc, "div", "storytext")
	if len(nodes) == 0 {
		return "", errors.New(`[-] There is no element match '<div id="storytext">'`)
	}
	plist := ElementsByTagName(nodes[0], "p")
	for _, v := range plist {
		if v.FirstChild == nil {
			continue
		} else if v.FirstChild.Data == "b" {
			body += "** "
			blist := ElementsByTagName(v, "b")
			for _, b := range blist {
				_b := b.FirstChild
				if _b != nil && _b.Data != "" {
					body += b.FirstChild.Data
				}
			}
			body += " **  \n"
		} else {
			body += v.FirstChild.Data + "  \n"
		}
	}
	body = strings.ReplaceAll(body, "**   **  \n", "")
	body = strings.ReplaceAll(body, "br  \n", "")
	return body, nil
}
