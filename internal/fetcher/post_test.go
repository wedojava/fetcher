package fetcher

import (
	"fmt"
	"testing"
	"time"
)

func TestSetAndSavePost(t *testing.T) {
	// p := PostFactory("https://www.dwnews.com/经济/60203253")
	p := PostFactory("https://www.dwnews.com/经济/60203034") // The wrong one
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		t.Errorf("GetRawAndDoC error: %v", err)
	}
	p.Raw, p.DOC = raw, doc
	if err := p.SetPost(); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	if err := p.SavePost(); err != nil {
		t.Errorf("test SetPost err: %v", doc)
	}
	fmt.Println(p.Title)
	fmt.Println(p.Body)
}
