package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/wedojava/gears"
)

// Posts is the fetcher can return many results.
type Posts map[string]*ThePost

type ThePost struct {
	Site   string
	Domain string
	URL    string
	Title  string
	Body   string
	Date   string
}

type Paragraph struct {
	Type    string
	Content string
}

func Fetch(url string) (*ThePost, error) {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		log.Fatal(err)
	}
	domain := gears.HttpGetDomain(url)
	site := gears.HttpGetSiteViaTwitterJS(rawBody)
	title := gears.HttpGetTitleViaTwitterJS(rawBody)
	// get contents
	body, err := FmtBodyDwnews(rawBody)
	if err != nil {
		log.Fatal(err)
	}
	date := gears.HttpGetDateViaMeta(rawBody)

	post := ThePost{site, domain, url, title, body, date}

	return &post, nil
}

func FetchUrls(url string) []string {
	rawBody, err := gears.HttpGetBody(url, 10)
	if err != nil {
		log.Fatal(err)
	}
	var ret_lst []string
	var reLink = regexp.MustCompile(`(?m)<a\shref\s?=\s?"(?P<href>/.{2}/\d{8}/.+?)".*?>`)
	for _, v := range reLink.FindAllStringSubmatch(rawBody, -1) {
		ret_lst = append(ret_lst, v[1])
	}
	ret_lst = gears.StrSliceDeDupl(ret_lst)

	return ret_lst
}

// FmtBodyDwnews focus on dwnews, it can extract raw body string via regexp and then, unmarshal it and format the news body to markdowned string.
func FmtBodyDwnews(rawBody string) (string, error) {
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
		jsTxtBody = strings.ReplaceAll(jsTxtBody, "],[", ",")
		jsTxtBody = jsTxtBody[:len(jsTxtBody)-1] + "]" // now body json data prepared done.
		// Unmarshal the json data
		var paragraph []Paragraph
		err := json.Unmarshal([]byte(jsTxtBody), &paragraph)
		if err != nil {
			return "", fmt.Errorf("[-] fetcher.FmtBodyDwnews()>Unmarshal() Error: %q", err)
		}
		for _, p := range paragraph {
			if p.Type == "boldText" {
				body += "**" + p.Content + "**  \n"
			} else {
				body += p.Content + "  \n"
			}

		}
	}

	return body, nil
}
