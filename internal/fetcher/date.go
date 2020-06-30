package fetcher

import (
	"fmt"
	"path/filepath"
	"regexp"

	"golang.org/x/net/html"
)

type dateItem struct {
	tagName             string
	positioningAttrName string
	positioningAttrVal  string
	dateAttrName        string
	dateAttrVal         string
	firstChildData      string
}

// TODO: return value and set same value to object at one function is redundancy
func soleDate(doc *html.Node, d *dateItem) (date string, err error) {
	type bailout struct{}
	defer func() {
		switch p := recover(); p {
		case nil:
			// no panic
		case bailout{}:
			// "expected" panic
			err = fmt.Errorf("multiple date elements")
		default:
			panic(p)
		}
	}()
	// Bail out of recursion if we find more than one non-empty date.
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == d.tagName {
			yes := false
			for _, a := range n.Attr {
				// TODO: compare items implement via map and deepcompare
				// yeap, positionAttr is blank means:
				// no need attr to help it to positioning where is date value
				// just tag name is enough to positioning the date value
				if d.positioningAttrName == "" && d.positioningAttrVal == "" {
					yes = true
				} else if a.Key == d.positioningAttrName && a.Val == d.positioningAttrVal {
					yes = true
				}
				if yes && d.dateAttrName != "" && a.Key == d.dateAttrName {
					date = a.Val
					d.dateAttrVal = a.Val
				} else if yes && d.dateAttrName == "" { // dateAttrName == "" means datevalue at firstChild
					date = n.FirstChild.Data
					d.dateAttrVal = date
				}
			}
		}
	}, nil)
	if date == "" {
		return "", fmt.Errorf("no date element")
	}
	return
}

func (p *Post) DwnewsDateInMeta() error {
	d := &dateItem{
		tagName:             "meta",
		positioningAttrName: "name",
		positioningAttrVal:  "parsely-pub-date",
		dateAttrName:        "content",
	}
	date, err := soleDate(p.DOC, d)
	if err != nil {
		return err
	}
	p.Date = date
	return nil
}

func (p *Post) VoaDateInNode() error {
	d := &dateItem{
		tagName:      "time",
		dateAttrName: "datetime",
	}
	date, err := soleDate(p.DOC, d)
	if err != nil {
		return err
	}
	p.Date = date
	return nil
}

func (p *Post) RfaDateInScript() error {
	d := &dateItem{
		tagName:             "script",
		positioningAttrName: "type",
		positioningAttrVal:  "application/ld+json",
	}
	date, err := soleDate(p.DOC, d)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`"date\w*?":\s*?"(.*?)"`)
	rs := re.FindAllStringSubmatch(date, -1)
	p.Date = rs[0][1] // dateModified -> rs[0][1], datePublished -> rs[1][1]
	return nil
}

func (p *Post) SetDate() error {
	switch p.Domain {
	case "www.boxun.com":
		a := filepath.Base(p.URL.Path)
		p.Date = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", a[:4], a[4:6], a[6:8], a[8:10], a[10:12], "00")
	case "www.dwnews.com":
		if err := p.DwnewsDateInMeta(); err != nil {
			return err
		}
	case "www.voachinese.com":
		if err := p.VoaDateInNode(); err != nil {
			return err
		}
	case "www.rfa.org":
		if err := p.RfaDateInScript(); err != nil {
			return err
		}
	}
	fmt.Println(p.Date) // print for test
	return nil
}
