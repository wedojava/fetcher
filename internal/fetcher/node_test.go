package fetcher

import (
	"fmt"
	"testing"
	"time"
)

func TestElementsByTagAndClass(t *testing.T) {
	p := PostFactory("https://www.voachinese.com/a/pandemic-drives-digital-innovations-in-u-s-presidential-race-20200701/5484814.html")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetDOC error: %v", err)
	}
	p.DOC = doc
	p.Raw = raw
	tc := ElementsByTagAndClass(doc, "div", "wsw")
	plist := ElementsByTagName(tc[0], "p")
	for _, v := range plist {
		fmt.Println(v.FirstChild.Data)
	}
}
