package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func main() {
	// ExampleDelRoutine()
	// DelRoutine("test", 3)
	// Output:
	// DelRoutine will delete: [04.23][1234H]test.txt
	// DelRoutine will delete: [04.24][1234H]test.txt
	// DelRoutine will delete: [04.25][1234H]test.txt
	// DelRoutine will delete: [04.26][1234H]test.txt
	// DelRoutine will delete: [04.26][2334H]test.txt
	fmt.Print("\n\n\n")
	fmt.Println("#========================================================#")
	fmt.Println("#======================= 新闻提取 =======================#")
	fmt.Println("#========================================================#")
	fmt.Println("#            呵~ 愚蠢的人类~ 你想要做些什么?!            #")
	fmt.Println("#========================================================#")
	fmt.Println("#                                                        #")
	fmt.Println("# [?] 不停的按回车可以退出此程序                         #")
	fmt.Println("#                                                        #")
	fmt.Println("# [1] 输入网页地址, 从网页直接提取                       #")
	fmt.Println("#                                                        #")
	fmt.Println("# [2] 按程序计划执行任务,目前只针对 多维新闻 展开        #")
	fmt.Println("#                                                        #")
	fmt.Println("#========================================================#")
	fmt.Printf("\n-> [?] 请输入选项序号: ")
	op := gears.GetInput()
	if strings.Compare("1", op) == 0 {
		FetchFromInput()
	} else if strings.Compare("2", op) == 0 {
		ServiceDwNews()
	} else {
		fmt.Printf("\nBye!\n\n")
		return
	}

}

// FetchFromInput fetch and save content to file by url input from terminal.
func FetchFromInput() {
	for {
		fmt.Printf("\n-> [!] 输入网页地址并回车(连续回车退出程序)：\n")
		url := gears.GetInput()
		if strings.Contains(url, "http") {
			SaveOne(url)
		} else {
			fmt.Printf("\nBye!\n\n")
			return
		}
	}
}

func ServiceDwNews() {
	var urlsNow, urlsBefore []string
	for {
		// 1. get url list from domain
		urlsNow = fetcher.FetchUrls("https://www.dwnews.com")
		// 2. compare urls, get diff urls between 2 lists then update urlsBefore and save.
		diff := gears.StrSliceDiff(urlsNow, urlsBefore)
		urlsBefore = urlsNow
		if len(diff) > 0 {
			for _, v := range diff {
				SaveOne("https://www.dwnews.com" + v)
			}
		}
		// TODO TO BE DISCUSSED: remove files that not contain in the pointed page.
		// Remove files 3 days ago
		DelRoutine("www.dwnews.com", 3)
		// all action above loop every 5 min.
		time.Sleep(5 * time.Minute)
		// *Optional. if the site folder is not exist or empty, means it's new action, so, the loop will action after first init files save.

	}
}

// SaveOne fetch content from url and save it if it not exist.
func SaveOne(url string) {
	f, _ := fetcher.Fetch(url)
	t, err := time.Parse(time.RFC3339, f.Date)
	if err != nil {
		fmt.Printf("\n[-] SaveOne()>time.Parse() error.\n%v\n", err)
	}
	filename := fmt.Sprintf("[%02d.%02d][%02d%02dH]%s%s", t.Month(), t.Day(), t.Hour(), t.Minute(), f.Title, ".md")
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(f.Domain)
	savePath := filepath.Join(f.Domain, filename)
	if !gears.Exists(savePath) {
		err = ioutil.WriteFile(filepath.Join(f.Site, filename), []byte(f.Body), 0644)
		if err != nil {
			fmt.Printf("\n[-] SaveOne()>WriteFile() error.\n%v\n", err)
		}
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
