package fetcher

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Fetcher struct {
	Entrance *url.URL
	Links    []string
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
func BreadthFirst(f func(item string), worklist []string) {
	for _, item := range worklist {
		f(item)
	}
}

func Crawl(_url string) {
	f := FetcherFactory(_url)
	log.Printf("[*] Deal with: [%s]\n", _url)
	log.Println("[*] Fetch links ...")
	if err := f.SetLinks(); err != nil { // f.Links update to the _url website is.
		log.Println(err)
		// if links cannot fetch sleep 1 minute then continue
		time.Sleep(1 * time.Minute)
		// continue // only useful by goroutine
		return
	}

	// GetNews then compare via md5 and Save or Rewrite news exist
	log.Println("[*] Get news ...")
	for _, link := range f.Links {
		post := PostFactory(link)
		if err := post.TreatPost(); err != nil {
			errMsg := "[-] TreatPost error occur on: " + link
			log.Println(errMsg)
			log.Println(err)
			ErrLog(errMsg + " " + err.Error())
		}
	}
	// Remove files 3 days ago
	DelRoutine(filepath.Join("wwwroot", f.Entrance.Hostname()), 3)
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
