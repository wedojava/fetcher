package fetcher

import (
	"fmt"
	"testing"
	"time"
)

func TestSetBody(t *testing.T) {
	// prepare
	// p := PostFactory("https://www.dwnews.com/%E5%85%A8%E7%90%83/60202451") // The right one
	// p := PostFactory("https://www.dwnews.com/%E5%8F%B0%E6%B9%BE/60202352") // The wrong one
	// p := PostFactory("https://www.rfa.org/mandarin/yataibaodao/shehui/hj-07022020095655.html")
	p := PostFactory("https://www.dwnews.com/经济/60203253") // The wrong one
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)

	if err != nil {
		t.Errorf("GetDOC error: %v", err)
	}
	p.DOC = doc
	p.Raw = raw
	// test
	if err := p.SetBody(); err != nil {
		t.Errorf("Error occur while SetBody() invoked: %v", err)
	}
	fmt.Println(p.Body)
}

func TestDwnews(t *testing.T) {
	// p := PostFactory("https://www.dwnews.com/%E5%85%A8%E7%90%83/60202451") // The right one
	// p := PostFactory("https://www.dwnews.com/%E5%8F%B0%E6%B9%BE/60202352") // The wrong one
	p := PostFactory("https://www.dwnews.com/全球/60203234")
	// p := PostFactory("https://www.dwnews.com/经济/60203253")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.DOC = doc
	p.Raw = raw
	body, err := Dwnews(p)
	if err != nil {
		t.Errorf("Dwnews error: %v", err)
	}
	fmt.Println(body)
}

func TestRfa(t *testing.T) {
	p := PostFactory("https://www.rfa.org/mandarin/yataibaodao/junshiwaijiao/jt-07022020105416.html")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	body, err := Rfa(p)
	if err != nil {
		t.Errorf("Voa error: %v", err)
	}
	fmt.Println(body)
}

func TestVoa(t *testing.T) {
	p := PostFactory("https://www.voachinese.com/a/controversial-national-security-law-enforced-in-hong-kong-despite-strong-opposition-from-us-and-hk-20200701/5484605.html")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	body, err := Voa(p)
	if err != nil {
		t.Errorf("Voa error: %v", err)
	}
	fmt.Println(body)
}

func TestBoxun(t *testing.T) {
	p := PostFactory("https://www.boxun.com/news/gb/china/2020/07/202007021503.shtml")
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	body, err := Boxun(p)
	if err != nil {
		t.Errorf("Voa error: %v", err)
	}
	fmt.Println(body)
}
