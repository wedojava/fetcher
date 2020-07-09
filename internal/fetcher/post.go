package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/wedojava/fetcher/internal/fetcher/sites/boxun"
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

// TODO: interupte while any err occur, the better way is write it done but don't interupt
func (p *Post) SetPost() error {
	if err := p.PostInit(); err != nil {
		return err
	}
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

// TODO: use func init
func (p *Post) PostInit() error {
	raw, doc, err := GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.Raw, p.DOC = raw, doc
	return nil
}

// TODO: rename as SetPost if it can be rm
func (p *Post) TreatPost() error {
	if err := p.PostInit(); err != nil {
		return err
	}
	switch p.Domain {
	case "www.boxun.com":
		// Set date
		post := boxun.Post(*p)
		p.Date = boxun.SetDate(&post)
		// Set title
		if err := p.TreatTitle(boxun.Title); err != nil {
			return err
		}
	}
	return nil
}

func (p *Post) SavePost() error {
	folderPath := filepath.Join("wwwroot", p.Domain)
	gears.MakeDirAll(folderPath)
	if p.Filename == "" {
		return errors.New("SavePost need a filename, but got none.")
	}
	filepath := filepath.Join(folderPath, p.Filename)
	if p.Body == "" {
		p.Body = "[-] Fetch error on visit: " + p.URL.String()
	}
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
