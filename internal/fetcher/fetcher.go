package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Posts is the fetcher can return many results.
type Posts map[string]*ThePost

type ThePost struct {
	url   string
	title string
	body  string
}

type Paragraph struct {
	Type    string
	Content string
}

func Fetch(url string) (*ThePost, error) {
	// Request the HTML page
	raw, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "[-] Fetch()>Get() Error!")
	}
	rawBody, err := ioutil.ReadAll(raw.Body)
	defer raw.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "[-] Fetch()>ReadAll() Error!")
	}
	if raw.StatusCode != 200 {
		return nil, errors.Wrap(err, "[-] Fetch()>Get() Error! Message: Cannot open the url.")
	}

	// var reTitle = regexp.MustCompile(`(?m)<title(.*?){0,1}>(?P<title>.*?)</title>`)
	var reTitle = regexp.MustCompile(`(?m)<meta name="twitter:title" content="(?P<title>.*?)"`)
	title := reTitle.FindStringSubmatch(string(rawBody))[1]
	// get contents
	// fetch and make it to json fmt
	var tmp = "["
	var reContent = regexp.MustCompile(`"htmlTokens":\[\[(?P<contents>.*?)\]\]`)
	for _, v := range reContent.FindAllStringSubmatch(string(rawBody), -1) {
		//fmt.Printf("index: %d => %v\n", i, v[1])
		tmp += v[1] + ","
	}
	tmp = strings.ReplaceAll(tmp, "],[", ",")
	tmp = tmp[:len(tmp)-1] + "]" // now body json data prepared done.
	// Unmarshal the json data
	b := []byte(tmp)
	var paragraph []Paragraph
	err = json.Unmarshal(b, &paragraph)
	if err != nil {
		return nil, fmt.Errorf("[-] Fetch()>Unmarshal() Error: %q", err)
	}
	//fmt.Printf("%+v", paragraph)
	// splice contents
	var body string
	for _, p := range paragraph {
		if p.Type == "boldText" {
			body += "**" + p.Content + "**  \n"
		} else {
			body += p.Content + "  \n"
		}

	}

	post := ThePost{url, title, body}

	return &post, nil
}
