package fetcher

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

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

func (p *Post) BoxunDateInUrl() error {
	rawdate := filepath.Base(p.URL.String())
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
		fmt.Println("err date fetch from url: ", p.URL.String())
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
	p.Date = fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02dZ", Y, M, D, hh, mm, 0)
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
	// fmt.Println(p.Date) // print for test
	return nil
}
