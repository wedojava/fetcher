package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	fetcherBoxun "github.com/wedojava/fetcher/internal/fetcher/boxun"
	fetcher "github.com/wedojava/fetcher/internal/fetcher/dwnews"
	fetcherRfa "github.com/wedojava/fetcher/internal/fetcher/rfaorg"
	fetcherVoa "github.com/wedojava/fetcher/internal/fetcher/voachinese"
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
	fmt.Println("# [?] Ctrl+c 或 不停的按回车可以退出此程序               #")
	fmt.Println("#                                                        #")
	fmt.Println("# [1] 按程序计划执行任务                                 #")
	fmt.Println("#                                                        #")
	fmt.Println("#========================================================#")
	fmt.Printf("\n-> [?] 请输入选项序号: ")
	op := gears.GetInput()
	if strings.Compare("1", op) == 0 {
		for {
			ServiceDwNews()
			ServiceRfa()
			ServiceVoa()
			ServiceBoxun()
			time.Sleep(5 * time.Minute)
		}
	} else if strings.Compare("2", op) == 0 {
		FetchFromInput()
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
			SaveOneDwnew(url)
		} else {
			fmt.Printf("\nBye!\n\n")
			return
		}
	}
}

func Service(srv func(start string)) {
	if srv != nil {
		srv("https://www.boxun.com/rolling.shtml")
	}
}

func ServiceBoxun() {
	var urlsNow, urlsBefore []string
	// for {
	// 1. get url list from domain
	urlsNow = fetcherBoxun.FetchBoxunUrls("https://boxun.com/rolling.shtml")
	// 2. compare urls, get diff urls between 2 lists then update urlsBefore and save.
	diff := gears.StrSliceDiff(urlsNow, urlsBefore)
	urlsBefore = urlsNow
	if len(diff) > 0 {
		for _, v := range diff {
			SaveOneBoxun(v)
		}
	}
	// Remove files 3 days ago
	DelRoutine(filepath.Join("wwwroot", "www.boxun.com"), 3)
	// }
}

func SaveOneBoxun(url string) {
	f, err := fetcherBoxun.FetchBoxun(url)
	if err != nil {
		fmt.Printf("\n[-] SaveOneBoxun()>FetchBoxun(%s) error occur:\n[-] %v", url, err)
		return
	}
	t, err := time.Parse(time.RFC3339, f.Date)
	if err != nil {
		fmt.Printf("\n[-] SaveOneBoxun()>time.Parse() error.\n%v\n", err)
		return
	}
	newsTitle := fmt.Sprintf("[%02d.%02d][%02d%02dH]%s", t.Month(), t.Day(), t.Hour(), t.Minute(), f.Title)
	filename := newsTitle + ".txt"
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(filepath.Join("wwwroot", f.Domain))
	savePath := filepath.Join("wwwroot", f.Domain, filename)
	if !gears.Exists(savePath) {
		err = ioutil.WriteFile(savePath, []byte("#"+newsTitle+"\n\n"+f.Body), 0644)
		if err != nil {
			fmt.Printf("\n[-] SaveOneBoxun()>WriteFile() error.\n%v\n", err)
		}
	}
}

func ServiceVoa() {
	var urlsNow, urlsBefore []string
	// for {
	// 1. get url list from domain
	urlsNow = fetcherVoa.FetchVoaUrls("https://www.voachinese.com")
	// 2. compare urls, get diff urls between 2 lists then update urlsBefore and save.
	diff := gears.StrSliceDiff(urlsNow, urlsBefore)
	urlsBefore = urlsNow
	if len(diff) > 0 {
		for _, v := range diff {
			SaveOneVoa(v)
		}
	}
	// Remove files 3 days ago
	DelRoutine(filepath.Join("wwwroot", "www.voachinese.com"), 3)
	// all action above loop every 5 min.
	// time.Sleep(5 * time.Minute)
	// *Optional. if the site folder is not exist or empty, means it's new action, so, the loop will action after first init files save.

	// }
}

func SaveOneVoa(url string) {
	f, err := fetcherVoa.FetchVoa(url)
	if err != nil {
		fmt.Printf("\n[-] SaveOneVoa()>FetchVoa(%s) error occur:\n[-] %v", url, err)
		return
	}
	t, err := time.Parse(time.RFC3339, f.Date)
	if err != nil {
		fmt.Printf("\n[-] SaveOneVoa()>time.Parse() error.\n%v\n", err)
		return
	}
	newsTitle := fmt.Sprintf("[%02d.%02d][%02d%02dH]%s", t.Month(), t.Day(), t.Hour(), t.Minute(), f.Title)
	filename := newsTitle + ".txt"
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(filepath.Join("wwwroot", f.Domain))
	savePath := filepath.Join("wwwroot", f.Domain, filename)
	if !gears.Exists(savePath) {
		err = ioutil.WriteFile(savePath, []byte("#"+newsTitle+"\n\n"+f.Body), 0644)
		if err != nil {
			fmt.Printf("\n[-] SaveOneVoa()>WriteFile() error.\n%v\n", err)
		}
	}
}

func ServiceRfa() {
	var urlsNow, urlsBefore []string
	// for {
	// 1. get url list from domain
	urlsNow = fetcherRfa.FetchRfaUrls("https://www.rfa.org/mandarin")
	// 2. compare urls, get diff urls between 2 lists then update urlsBefore and save.
	diff := gears.StrSliceDiff(urlsNow, urlsBefore)
	urlsBefore = urlsNow
	if len(diff) > 0 {
		for _, v := range diff {
			SaveOneRfa(v)
		}
	}
	// Remove files 3 days ago
	DelRoutine(filepath.Join("wwwroot", "www.rfa.org"), 3)
	// all action above loop every 5 min.
	// time.Sleep(5 * time.Minute)
	// *Optional. if the site folder is not exist or empty, means it's new action, so, the loop will action after first init files save.

	// }
}

func SaveOneRfa(url string) {
	f, err := fetcherRfa.FetchRfa(url)
	if err != nil {
		fmt.Printf("\n[-] SaveOneRfa()>FetchRfa(%s) error occur:\n[-] %v", url, err)
		return
	}
	t, err := time.Parse(time.RFC3339, f.Date)
	if err != nil {
		fmt.Printf("\n[-] SaveOneRfa()>time.Parse() error.\n%v\n", err)
		return
	}
	newsTitle := fmt.Sprintf("[%02d.%02d][%02d%02dH]%s", t.Month(), t.Day(), t.Hour(), t.Minute(), f.Title)
	filename := newsTitle + ".txt"
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(filepath.Join("wwwroot", f.Domain))
	savePath := filepath.Join("wwwroot", f.Domain, filename)
	if !gears.Exists(savePath) {
		err = ioutil.WriteFile(savePath, []byte("#"+newsTitle+"\n\n"+f.Body), 0644)
		if err != nil {
			fmt.Printf("\n[-] SaveOneRfa()>WriteFile() error.\n%v\n", err)
		}
	}
}

func ServiceDwNews() {
	var urlsNow, urlsBefore []string
	// for {
	// 1. get url list from domain
	urlsNow = fetcher.FetchDwnewsUrls("https://www.dwnews.com")
	// 2. compare urls, get diff urls between 2 lists then update urlsBefore and save.
	diff := gears.StrSliceDiff(urlsNow, urlsBefore)
	urlsBefore = urlsNow
	if len(diff) > 0 {
		for _, v := range diff {
			SaveOneDwnew(v)
		}
	}
	// Remove files 3 days ago
	DelRoutine(filepath.Join("wwwroot", "www.dwnews.com"), 3)
	// all action above loop every 5 min.
	// time.Sleep(5 * time.Minute)
	// *Optional. if the site folder is not exist or empty, means it's new action, so, the loop will action after first init files save.

	// }
}

// SaveOne fetch content from url and save it if it not exist.
func SaveOneDwnew(url string) {
	f, err := fetcher.FetchDwnews(url)
	if err != nil {
		fmt.Printf("\n[-] SaveOneDwnew()>FetchDwnews(%s) error occur:\n[-] %v", url, err)
		return
	}
	t, err := time.Parse(time.RFC3339, f.Date)
	if err != nil {
		fmt.Printf("\n[-] SaveOneDwnew()>time.Parse() error.\n%v\n", err)
		return
	}
	newsTitle := fmt.Sprintf("[%02d.%02d][%02d%02dH]%s", t.Month(), t.Day(), t.Hour(), t.Minute(), f.Title)
	filename := newsTitle + ".txt"
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(filepath.Join("wwwroot", f.Domain))
	savePath := filepath.Join("wwwroot", f.Domain, filename)
	if !gears.Exists(savePath) {
		err = ioutil.WriteFile(savePath, []byte("#"+newsTitle+"\n\n"+f.Body), 0644)
		if err != nil {
			fmt.Printf("\n[-] SaveOneDwnew()>WriteFile() error.\n%v\n", err)
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
