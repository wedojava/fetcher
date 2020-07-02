package fetcher

import (
	"fmt"
	"testing"
	"time"
)

func TestSetBody(t *testing.T) {
	// prepare
	// p := PostFactory("https://www.dwnews.com/%E5%85%A8%E7%90%83/60202451") // The right one
	p := PostFactory("https://www.dwnews.com/%E5%8F%B0%E6%B9%BE/60202352") // The wrong one
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
	p := PostFactory("https://www.dwnews.com/%E5%85%A8%E7%90%83/60202451") // The right one
	// p := PostFactory("https://www.dwnews.com/%E5%8F%B0%E6%B9%BE/60202352") // The wrong one
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
