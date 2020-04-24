package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func main() {
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
	// fmt.Println("# [2] 从当前目录下 `list.txt` 文件中获得下载列表      #")
	// fmt.Println("#                                                        #")
	fmt.Println("#========================================================#")
	fmt.Printf("\n-> [?] 请输入选项序号: ")
	reader := bufio.NewReader(os.Stdin)
	op, _ := reader.ReadString('\n')
	// convert CRLF to LF
	op = strings.Replace(op, "\n", "", -1)
	op = strings.Replace(op, "\r", "", -1)
	if strings.Compare("1", op) == 0 {
		for {
			fmt.Printf("\n-> [!] 输入网页地址并回车(连续回车退出程序)：\n")
			url, _ := reader.ReadString('\n')
			// convert CRLF to LF
			url = strings.Replace(url, "\n", "", -1)
			url = strings.Replace(url, "\r", "", -1)
			// fmt.Scanf("%s", &url)
			if strings.Contains(url, "http") {
				SaveOne(url)
			} else {
				fmt.Printf("\nBye!\n\n")
				return
			}
		}
	} else if strings.Compare("2", op) == 0 {
		fmt.Printf("\n-> [*] 批量内容提取开始......\n")
	} else {
		fmt.Printf("\nBye!\n\n")
		return
	}

}

func SaveOne(url string) {
	f, _ := fetcher.Fetch(url)
	// Save Body to file named title in folder twitter site content
	gears.MakeDirAll(f.Site)
	err := ioutil.WriteFile(filepath.Join(f.Site, f.Title+".md"), []byte(f.Body), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
