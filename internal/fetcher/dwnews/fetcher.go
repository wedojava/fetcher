package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/wedojava/fetcher/internal/fetcher"
	"github.com/wedojava/gears"
)

func FetchDwnews(url string) (*fetcher.ThePost, error) {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		return nil, err
		// log.Fatal(err)
	}
	domain := gears.HttpGetDomain(url)
	site := gears.HttpGetSiteViaTwitterJS(rawBody)
	title := ThisGetTitle(rawBody)
	// get contents
	body, err := FmtBodyDwnews(rawBody)
	if err != nil {
		fmt.Println(err)
		return nil, err
		// log.Fatal(err)
	}
	date := gears.HttpGetDateViaMeta(rawBody)

	post := fetcher.ThePost{Site: site, Domain: domain, URL: url, Title: title, Body: body, Date: date}

	return &post, nil
}

func ThisGetTitle(rawBody string) string {
	var a = regexp.MustCompile(`(?m)<h1.*?id="articleTitle".*?>(.*?)</h1>`)
	rt := a.FindStringSubmatch(rawBody)
	if rt != nil {
		return rt[1]

	} else {
		return ""

	}
}

func FetchDwnewsUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<a\shref\s?=\s?"(?P<href>/.{2}/\d{8}/.+?)".*?>`)
	lst := reLink.FindAllStringSubmatch(rawBody, -1)
	if lst == nil {
		fmt.Printf("\n[-] fetcher.FetchDwnewsUrls(%s) regex matched nothing.\n", url)
		return nil
	} else {
		for _, v := range reLink.FindAllStringSubmatch(rawBody, -1) {
			ret_lst = append(ret_lst, "https://www.dwnews.com"+v[1])
		}
		ret_lst = gears.StrSliceDeDupl(ret_lst)
	}

	return ret_lst
}

// FmtBodyDwnews focus on dwnews, it can extract raw body string via regexp and then, unmarshal it and format the news body to markdowned string.
func FmtBodyDwnews(rawBody string) (string, error) {
	if rawBody == "" {
		return "", errors.New("[-] FmtBodyDwnews() parameter is nil!")
	}
	// extract and make it to json fmt
	var jsTxtBody = "["
	var body string // splice contents
	var reContent = regexp.MustCompile(`"htmlTokens":\[\[(?P<contents>.*?)\]\]`)
	for _, v := range reContent.FindAllStringSubmatch(rawBody, -1) {
		jsTxtBody += v[1] + ","
	}
	if jsTxtBody == "[" { // this means jsTxtBody got northing, so it may be pic news.
		reContent = regexp.MustCompile(`"\d{7}":{"caption":"(?P<title>.*?)"`)
		for _, v := range reContent.FindAllStringSubmatch(rawBody, -1) {
			body += v[1] + "  \n"
		}
	} else {
		var reSummary = regexp.MustCompile(`"blockType":"summary","summary":\[(".*?")\]}`)
		rs := reSummary.FindStringSubmatch(rawBody)
		if rs != nil {
			jsTxtBody = fmt.Sprintf("[{\"type\":\"summary\",\"content\":%s}],%s", rs[1], jsTxtBody)
		}
		jsTxtBody = strings.ReplaceAll(jsTxtBody, "],[", ",")
		jsTxtBody = jsTxtBody[:len(jsTxtBody)-1] + "]" // now body json data prepared done.
		// Unmarshal the json data
		var paragraph []fetcher.Paragraph
		err := json.Unmarshal([]byte(jsTxtBody), &paragraph)
		if err != nil {
			return "", fmt.Errorf("[-] fetcher.FmtBodyDwnews()>Unmarshal() Error: %q", err)
		}
		for _, p := range paragraph {
			if p.Type == "boldText" {
				body += "**" + p.Content + "**  \n"
			} else if p.Type == "summary" {
				body += "\n> " + p.Content + "  \n  \n"
			} else {
				body += p.Content + "  \n"
			}

		}
	}

	return body, nil
}
