package htmldoc

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestElementsByTagAndClass(t *testing.T) {
	u, err := url.Parse("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/jt-07022020105416.html")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndClass(doc, "div", "wsw")
	plist := ElementsByTagName(tc[0], "p")
	for _, v := range plist {
		fmt.Println(v.FirstChild.Data)
	}
}

func TestElementsByTagAndId(t *testing.T) {
	u, err := url.Parse("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/jt-07022020105416.html")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := ElementsByTagAndId(doc, "div", "storytext")
	plist := ElementsByTagName(tc[0], "p")
	for _, v := range plist {
		if v.FirstChild != nil {
			if v.FirstChild.Data == "b" {
				blist := ElementsByTagName(v, "b")
				fmt.Print("**")
				for _, b := range blist {
					fmt.Print(b.FirstChild.Data)
				}
				fmt.Print("**\n")
				// fmt.Println("**" + v.FirstChild.FirstChild.Data + "**")
			} else {
				fmt.Println(v.FirstChild.Data)
			}
		}
	}
}

func TestMetaByName(t *testing.T) {
	u, err := url.Parse("https://www.dwnews.com/全球/60203304")
	if err != nil {
		t.Errorf("url Parse err: %v", err)
	}
	_, doc, err := GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoc err: %v", err)
	}
	tc := MetasByName(doc, "parsely-pub-date")
	rt := []string{}
	for _, n := range tc {
		for _, a := range n.Attr {
			if a.Key == "content" {
				rt = append(rt, a.Val)
			}
		}
	}
	want := "2020-07-09T18:04:00+08:00"
	if want != rt[0] {
		t.Errorf("want: %v, got: %v", want, rt[0])
	}
	fmt.Println(rt[0])
}
