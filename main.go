package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func main() {
	f, _ := fetcher.Fetch("https://www.dwnews.com/中国/60176170")
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(f.Site)
	err := ioutil.WriteFile(filepath.Join(f.Site, f.Title+".md"), []byte(f.Body), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
