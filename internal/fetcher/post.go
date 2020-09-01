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
	"github.com/wedojava/fetcher/internal/fetcher/sites/dwnews"
	"github.com/wedojava/fetcher/internal/fetcher/sites/ltn"
	"github.com/wedojava/fetcher/internal/fetcher/sites/rfa"
	"github.com/wedojava/fetcher/internal/fetcher/sites/voachinese"
	"github.com/wedojava/fetcher/internal/fetcher/sites/zaobao"
	"github.com/wedojava/fetcher/internal/htmldoc"
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

// TODO: use func init
func (p *Post) PostInit() error {
	raw, doc, err := htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	if err != nil {
		return err
	}
	p.Raw, p.DOC = raw, doc
	return nil
}

// TreatPost get post things and set to `p` then save it.
func (p *Post) TreatPost() error {
	// Init post
	if err := p.PostInit(); err != nil {
		return err
	}
	// Set post
	switch p.Domain {
	case "www.boxun.com":
		post := boxun.Post(*p)
		if err := boxun.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.dwnews.com":
		post := dwnews.Post(*p)
		if err := dwnews.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.voachinese.com":
		post := voachinese.Post(*p)
		if err := voachinese.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.rfa.org":
		post := rfa.Post(*p)
		if err := rfa.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "www.zaobao.com":
		post := zaobao.Post(*p)
		if err := zaobao.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	case "news.ltn.com.tw":
		post := ltn.Post(*p)
		if err := ltn.SetPost(&post); err != nil {
			return err
		}
		*p = Post(post)
	}
	// Save post to file
	if err := p.setFilename(); err != nil {
		return err
	}
	if err := p.savePost(); err != nil {
		return err
	}
	return nil
}

func (p *Post) savePost() error {
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
	return nil
}

func (p *Post) setFilename() error {
	t, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		return err
	}
	p.Filename = fmt.Sprintf("[%02d.%02d][%02d%02dH]%s.txt", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	return nil
}
