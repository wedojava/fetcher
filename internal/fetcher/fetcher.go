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
	Site  string
	URL   string
	Title string
	Body  string
}

type Paragraph struct {
	Type    string
	Content string
}

func Fetch(url string) (*ThePost, error) {
	rawBody, err := gears.HttpGetBody(url)
	if err != nil {
		log.Fatal(err)
	}
	site := gears.HttpGetSiteViaTwitterJS(rawBody)
	title := gears.HttpGetTitleViaTwitterJS(rawBody)
	// get contents
	body, err := FmtBodyDwnews(rawBody)
	if err != nil {
		log.Fatal(err)
	}

	post := ThePost{site, url, title, body}

	return &post, nil
}

// FmtBodyDwnews focus on dwnews, it can extract raw body string via regexp and then, unmarshal it and format the news body to markdowned string.
func FmtBodyDwnews(rawBody string) (string, error) {
	// extract and make it to json fmt
	var tmp = "["
	var reContent = regexp.MustCompile(`"htmlTokens":\[\[(?P<contents>.*?)\]\]`)
	for _, v := range reContent.FindAllStringSubmatch(rawBody, -1) {
		tmp += v[1] + ","
	}
	tmp = strings.ReplaceAll(tmp, "],[", ",")
	tmp = tmp[:len(tmp)-1] + "]" // now body json data prepared done.
	// Unmarshal the json data
	var paragraph []Paragraph
	err := json.Unmarshal([]byte(tmp), &paragraph)
	if err != nil {
		return "", fmt.Errorf("[-] GetBodyDwnews()>Unmarshal() Error: %q", err)
	}
	// splice contents
	var body string
	for _, p := range paragraph {
		if p.Type == "boldText" {
			body += "**" + p.Content + "**  \n"
		} else {
			body += p.Content + "  \n"
		}

	}
	return body, nil
}
