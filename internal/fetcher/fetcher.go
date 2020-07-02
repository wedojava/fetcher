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
	Entrance *url.URL
	Links    []string
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
	u, err := url.Parse(site)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Fetcher{
		Entrance: u,
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func breadthFirst(f func(item string), worklist []string) {
	for _, item := range worklist {
		// TODO: run in goroutine
		f(item)
	}
}

// TODO: Links can be managed by each object fetcher.
func crawl(_url string) {
	f := FetcherFactory(_url)
	for {
		log.Printf("[*] Deal with: [%s]\n", _url)
		log.Println("[*] Fetch links ...")
		if err := f.SetLinks(); err != nil { // f.Links update to the _url website is.
			log.Println(err)
			// if links cannot fetch sleep 1 minute then continue
			time.Sleep(1 * time.Minute)
			continue
		}

		// GetNews then compare via md5 and Save or Rewrite news exist
		log.Println("[*] Get news ...")
		for _, link := range f.Links {
			post := PostFactory(link)
			if err := post.SetPost(); err != nil {
				errMsg := "[-] SetPost error occur from: " + link
				log.Printf(errMsg)
				log.Println(err)
				ErrLog(errMsg + " " + err.Error())
			}
			if err := post.SavePost(); err != nil {
				errMsg := "[-] SavePost error occur from: " + link
				log.Printf(errMsg)
				log.Println(err)
				ErrLog(errMsg + " " + err.Error())
			}
		}
		// Remove files 3 days ago
		DelRoutine(filepath.Join("wwwroot", f.Entrance.Hostname()), 3)
		// hold on 5 minutes
		log.Println("Sleep a sec ...")
		time.Sleep(5 * time.Minute)
	}
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
	write.WriteString("[" + time.Now().Format(time.RFC3339) + "]--------------------------------------\n")
	write.WriteString(msg + "\n")
	write.Flush()
	return nil
}
