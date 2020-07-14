package rfa

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/wedojava/fetcher/internal/htmldoc"
	"github.com/wedojava/gears"
	"golang.org/x/net/html"
)

type Post struct {
	Domain   string
	URL      *url.URL
	DOC      *html.Node
	Raw      []byte
	Title    string
	Body     string
	Date     string
	Filename string
}

func SetPost(p *Post) error {
	if err := SetDate(p); err != nil {
		return err
	}
	if err := SetTitle(p); err != nil {
		return err
	}
	if err := SetBody(p); err != nil {
		return err
	}
	return nil
}

func SetDate(p *Post) error {
	doc := htmldoc.ElementsByTagAndType(p.DOC, "script", "application/ld+json")
	if doc == nil {
		return errors.New("[-] rfa SetDate err, cannot get target nodes.")
	}
	d := doc[0].FirstChild
	if d.Type != html.TextNode {
		return errors.New("[-] rfa SetDate err, target node have no text.")
	}
	raw := d.Data
	re := regexp.MustCompile(`"date\w*?":\s*?"(.*?)"`)
	rs := re.FindAllStringSubmatch(raw, -1)
	p.Date = rs[0][1] // dateModified -> rs[0][1], datePublished -> rs[1][1]
	return nil
}

func SetTitle(p *Post) error {
	n := htmldoc.ElementsByTag(p.DOC, "title")
	if n == nil {
		return fmt.Errorf("[-] there is no element <title>")
	}
	title := n[0].FirstChild.Data
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func SetBody(p *Post) error {
	if p.DOC == nil {
		return errors.New("[-] there is no DOC object to get and format.")
	}
	b, err := Rfa(p)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	h1 := fmt.Sprintf("# [%02d.%02d][%02d%02dH] %s", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	p.Body = fmt.Sprintf("%s\n\n%s", h1, b)
	return nil
}

func Rfa3(p *Post) (string, error) {
	if p.Raw == nil {
		return "", errors.New("\n[-] FmtBodyRfa() parameter is nil!\n")
	}
	var ps []string
	var body string
	var reContent = regexp.MustCompile(`(?m)<p.*?>(?P<content>.*?)</p>`)
	for _, v := range reContent.FindAllStringSubmatch(string(p.Raw), -1) {
		ps = append(ps, v[1])
	}
	if len(ps) == 0 {
		if regexp.MustCompile(`(?m)<video.*?>`).FindAllString(string(p.Raw), -1) != nil {
			return "", errors.New("\n[-] fetcher.FmtBodyRfa() Error: this is a video page.\n")
		}
		return "", errors.New("\n[-] fetcher.FmtBodyRfa() Error: regex matched nothing.\n")
	} else {
		for _, p := range ps {
			body += p + "  \n"
		}
	}

	return body, nil
}

func Rfa2(p *Post) (string, error) {
	r := htmldoc.ElementsByTagAndId2(p.Raw, "div", "storytext")
	// z := html.NewTokenizerFragment(bytes.NewReader(r), "p")
	// fmt.Println(string(r))
	z := html.NewTokenizer(bytes.NewReader(r))
	var b bytes.Buffer
	for {
		tt := z.Next()
		token := z.Token()
		if err := z.Err(); err == io.EOF {
			break
		}
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				break
			} else {
				return "", fmt.Errorf("Rf2 err occur: %v", z.Err())
			}
		case html.TextToken:
			b.Write(z.Text())
		case html.SelfClosingTagToken:
			if "br" == token.Data {
				// z.Next()
				// z.Next()
				// z.Next()
				// b.Write(z.Text())
				continue
			}
		case html.StartTagToken, html.EndTagToken:
			if "p" == token.Data {
				// z.Next()
				// fmt.Println(z.Next().String())
				// b.Write(z.Buffered())
			}
			// if "b" == token.Data {
			//         z.Next()
			//         z.Next()
			//         b.Write([]byte("** "))
			//         b.Write(z.Text())
			//         b.Write([]byte(" **  \n"))
			// }
		}
	}
	return b.String(), nil
}

func Rfa(p *Post) (string, error) {
	doc := p.DOC
	body := ""
	// Fetch content nodes
	nodes := htmldoc.ElementsByTagAndId(doc, "div", "storytext")
	if len(nodes) == 0 {
		return "", errors.New(`[-] There is no element match '<div id="storytext">'`)
	}
	plist := htmldoc.ElementsByTag(nodes[0], "p")
	if len(plist) == 1 {
		innerNodes := htmldoc.ElementsNext(plist[0])
		for _, in := range innerNodes {
			if in.Type == html.TextNode {
				body += in.Data + "  \n"
			}
		}
	} else {
		for _, v := range plist {
			if v.FirstChild == nil {
				continue
			}
			htmldoc.ElementsRmByTag(v, "br")
			fd := v.FirstChild.Data
			if fd == "iframe" || fd == "i" {
				continue
			}
			if fd == "b" {
				body += "** "
				blist := htmldoc.ElementsByTag(v, "b")
				for _, b := range blist {
					_b := b.FirstChild
					if _b != nil && _b.Data != "" {
						body += _b.Data
					}
				}
				body += " **  \n"
			} else {
				body += fd + "  \n"
			}
		}
	}

	body = strings.ReplaceAll(body, "**   **  \n", "")
	body = strings.ReplaceAll(body, "br  \n", "")
	return body, nil
}
