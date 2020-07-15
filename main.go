package main

import (
	"fmt"
	"log"
	"strconv"
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
	fmt.Println("# [?] Ctrl+c 或 不停的按回车可以退出此程序               #")
	fmt.Println("#                                                        #")
	fmt.Println("# [1] 按程序计划执行任务                                 #")
	fmt.Println("#                                                        #")
	fmt.Println("#========================================================#")
	fmt.Printf("\n-> [?] 请输入选项序号: ")
	op := gears.GetInput()
	year := strconv.Itoa(time.Now().Year())
	sites := []string{
		"https://www.dwnews.com",
		"https://www.voachinese.com",
		"https://www.voachinese.com/z/1739",
		"https://www.rfa.org/mandarin/",
		"https://www.rfa.org/mandarin/Xinwen/story_archive?year=" + year,
		"https://www.rfa.org/mandarin/yataibaodao/story_archive?year=" + year,
		"https://www.boxun.com/rolling.shtml",
	}
	if strings.Compare("1", op) == 0 {
		for {
			fetcher.BreadthFirst(fetcher.Crawl, sites)
			log.Println("Hold a sec ...")
			time.Sleep(5 * time.Minute)
		}
	} else if strings.Compare("2", op) == 0 {
		// FetchFromInput()
	} else {
		fmt.Printf("\nBye!\n\n")
		return
	}

}
