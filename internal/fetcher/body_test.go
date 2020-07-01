package fetcher

import (
	"fmt"
	"testing"
	"time"
)

func TestSetBody(t *testing.T) {
	// prepare
	p := PostFactory("https://www.boxun.com/news/gb/intl/2020/06/202006302339.shtml")
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
