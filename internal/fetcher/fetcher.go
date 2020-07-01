package fetcher

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/wedojava/gears"
	"golang.org/x/net/html"
)

type Fetcher struct {
	Entrance string
	Links    []string
	LinksNew []string
	LinksOld []string
	LinksErr []string
}

var originalHost string

// GetRawAndDoc can get html raw bytes and html.Node by rawurl.
func GetRawAndDoc(url *url.URL, retryTimeout time.Duration) ([]byte, *html.Node, error) {
	// To judge if there is a syntex error on url
	if originalHost == "" {
		originalHost = url.Host
	}
	if originalHost != url.Host {
		return nil, nil, fmt.Errorf("bad host of url: %s, expected: %s", url.Host, originalHost)
	}
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := http.Get(url.String())
		if err == nil { // success
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, nil, err
			}
			reader := bytes.NewBuffer(raw)
			doc, err := html.Parse(reader)
			return raw, doc, nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil, nil, nil
}

func FetcherFactory(site string) *Fetcher {
	return &Fetcher{
		Entrance: site,
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func breadthFirst(f func(item string), worklist []string) {
	for _, item := range worklist {
		f(item)
	}
}

func crawl(_url string) {
	f := FetcherFactory(_url)
	log.Printf("[*] Deal with: [%s]\n", _url)
	log.Println("[*] Fetch links ...")
	if err := f.SetLinks(); err != nil {
		log.Println(err)
	}
	// Set LinksNew
	f.LinksNew = gears.StrSliceDiff(f.Links, f.LinksOld)
	if len(f.LinksNew) == 0 { // there's no news need to be saved.
		return
	}
	// GetNews then compare via md5 and Save or Rewrite news exist
	log.Println("[*] Get news ...")
	for _, link := range f.LinksNew {
		post := PostFactory(link)
		if err := post.SetPost(); err != nil {
			f.LinksErr = append(f.LinksErr, link)
			errMsg := "[-] SetPost error occur from: " + link
			log.Printf(errMsg)
			log.Println(err)
			ErrLog(errMsg + " " + err.Error())
		}
		if err := post.SavePost(); err != nil {
			f.LinksErr = append(f.LinksErr, link)
			errMsg := "[-] SavePost error occur from: " + link
			log.Printf(errMsg)
			log.Println(err)
			ErrLog(errMsg + " " + err.Error())
		}
	}
	// Set LinksOld, if only success above, then set LinksOld = Links
	f.LinksOld = f.Links
	// Remove files 3 days ago
	u, err := url.Parse(f.Entrance)
	if err != nil {
		log.Println(err)
		ErrLog(err.Error())
	}
	DelRoutine(filepath.Join("wwwroot", u.Hostname()), 3)
}

// DelRoutine remove files in folder days ago
func DelRoutine(folder string, n int) error {
	if !gears.Exists(folder) {
		fmt.Printf("\n[-] DelRoutine() err: Folder(%s) does not exist.\n", folder)
		return nil
	}
	for i := 0; i < 3; i++ { // deal with file n+i days ago
		a := time.Now().AddDate(0, 0, -(n + i))
		b := fmt.Sprintf("[%02d.%02d]", a.Month(), a.Day())
		c, _ := gears.GetPrefixedFiles(folder, b)
		for _, f := range c {
			// fmt.Println("DelRoutine will delete: ", f)
			os.Remove(f)
		}

	}
	return nil
}

func ErrLog(msg string) error {
	filePath := "./errLog.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(msg)
	write.Flush()
	return nil
}
