package main

import (
	"fmt"

	"github.com/wedojava/fetcher/internal/fetcher"
)

func main() {
	f, _ := fetcher.Fetch("https://www.dwnews.com/中国/60176170")
	fmt.Printf("%v", f)
}
