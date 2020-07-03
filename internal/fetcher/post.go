package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"time"

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

type Paragraph struct {
	Type    string
	Content string
}

func PostFactory(rawurl string) *Post {
	url, err := url.Parse(rawurl)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Post{
		Domain: url.Hostname(),
		URL:    url,
	}
}

func (p *Post) SetPost() error {
	// set contents
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.Raw, p.DOC = raw, doc
	// set Date
	if err := p.SetDate(); err != nil {
		return err
	}
	// set Title
	if err := p.SetTitle(); err != nil {
		return err
	}
	// set Filename
	if err := p.SetFilename(); err != nil {
		return err
	}
	// set Body (get and format body)
	if err := p.SetBody(); err != nil {
		return err
	}
	return nil
}

func (p *Post) SavePost() error {
	folderPath := filepath.Join("wwwroot", p.Domain)
	gears.MakeDirAll(folderPath)
	filepath := filepath.Join(folderPath, p.Filename)
	err := ioutil.WriteFile(filepath, []byte(p.Body), 0644)
	if err != nil {
		return err
	}
	// if !gears.Exists(filepath) {
	//         err := ioutil.WriteFile(filepath, []byte(p.Body), 0644)
	//         if err != nil {
	//                 return err
	//         }
	// }
	return nil
}

func (p *Post) SetFilename() error {
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	p.Filename = fmt.Sprintf("[%02d.%02d][%02d%02dH]%s.txt", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	return nil
}
